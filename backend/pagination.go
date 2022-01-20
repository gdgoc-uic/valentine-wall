package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type PaginatedLinks struct {
	First    string  `json:"first"`
	Last     string  `json:"last"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}

type PaginatedResponse struct {
	Links      PaginatedLinks `json:"links"`
	Page       int64          `json:"page"`
	PerPage    int64          `json:"per_page"`
	PageCount  int64          `json:"page_count"`
	TotalCount int64          `json:"total_count"`
	Data       interface{}    `json:"data"`
}

type Paginator struct {
	Path      *url.URL
	Page      int64
	Limit     int64
	OrderKey  string
	Order     string
	TableName string
}

func (pg Paginator) Copy(path *url.URL, page int64, limit int64, order string) Paginator {
	return Paginator{
		Path:     path,
		OrderKey: pg.OrderKey,
		Page:     page,
		Limit:    limit,
		Order:    order,
	}
}

func (pg Paginator) Filters(queryB sq.SelectBuilder) sq.SelectBuilder {
	fmt.Printf("{order: %s, limit: %d, page: %d, OrderKey: %s}\n", pg.Order, pg.Limit, pg.Page, pg.OrderKey)
	return queryB.OrderBy(pg.OrderKey + " " + pg.Order).Limit(uint64(pg.Limit)).Offset(uint64(pg.Page-1) * uint64(pg.Limit))
}

func generatePaginateUrl(fromUrl *url.URL, page int64, limit int64, order string) *string {
	nextUrl := *fromUrl
	queries := nextUrl.Query()
	queries.Set("page", fmt.Sprintf("%d", page))
	queries.Set("limit", fmt.Sprintf("%d", limit))
	queries.Set("order", order)
	nextUrl.RawQuery = queries.Encode()
	nextLink := nextUrl.String()
	return &nextLink
}

func (pg Paginator) Load(db *sqlx.DB, tableName string, queryB sq.SelectBuilder, converter func(*sqlx.Rows) (interface{}, error)) (*PaginatedResponse, error) {
	commaCount := strings.Count(pg.Order, ",")
	if commaCount == 1 {
		splitted := strings.Split(pg.Order, ",")
		pg.OrderKey = splitted[0]
		pg.Order = splitted[1]
	} else if commaCount > 1 {
		return nil, &ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "ordering by specific column must be column name follow by a single comma and order (asc/desc)",
		}
	}

	count := int64(0)
	if countQuery, _, err := sq.Select("count(*)").From(tableName).ToSql(); err != nil {
		return nil, err
	} else if err := db.Get(&count, countQuery); err != nil {
		return nil, err
	}

	lastPageNumber := count / pg.Limit
	if count%pg.Limit != 0 {
		lastPageNumber++
	}

	currentPageNumber := pg.Page
	if lastPageNumber == 0 {
		lastPageNumber = 1
	}

	if currentPageNumber > lastPageNumber || currentPageNumber < 1 {
		return nil, &ResponseError{StatusCode: http.StatusNotFound}
	}

	dataSql, args, err := pg.Filters(queryB.From(tableName)).ToSql()
	if err != nil {
		return nil, err
	}

	results := []interface{}{}
	rows, err := db.Queryx(dataSql, args...)
	if err != nil {
		return nil, &ResponseError{StatusCode: http.StatusNotFound, WError: err}
	}

	for rows.Next() {
		gotData, err := converter(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, gotData)
	}

	resp := &PaginatedResponse{
		Links: PaginatedLinks{
			First: *generatePaginateUrl(pg.Path, 1, pg.Limit, pg.OrderKey+" "+pg.Order),
			Last:  *generatePaginateUrl(pg.Path, lastPageNumber, pg.Limit, pg.OrderKey+" "+pg.Order),
		},
		Page:       currentPageNumber,
		TotalCount: count,
		PerPage:    pg.Limit,
		PageCount:  lastPageNumber,
		Data:       results,
	}

	if currentPageNumber+1 <= lastPageNumber {
		resp.Links.Next = generatePaginateUrl(pg.Path, currentPageNumber+1, pg.Limit, pg.OrderKey+" "+pg.Order)
	}

	if currentPageNumber > 1 {
		resp.Links.Previous = generatePaginateUrl(pg.Path, currentPageNumber-1, pg.Limit, pg.OrderKey+" "+pg.Order)
	}

	return resp, nil
}

func pagination(pg Paginator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			pageNumber := 1
			limitCount := 10
			order := "asc"

			if rawOrder := r.URL.Query().Get("order"); len(rawOrder) != 0 {
				splitted := strings.Split(rawOrder, ",")
				lastIndex := len(splitted) - 1
				if splitted[lastIndex] != "desc" && splitted[lastIndex] != "asc" {
					return &ResponseError{
						StatusCode: http.StatusBadRequest,
						Message:    "order invalid value",
					}
				}
				order = rawOrder
			}

			if rawLimitCount := r.URL.Query().Get("limit"); len(rawLimitCount) != 0 {
				var err error
				if limitCount, err = strconv.Atoi(rawLimitCount); err != nil {
					return &ResponseError{
						StatusCode: http.StatusBadRequest,
						Message:    "limit invalid value",
					}
				}
			}

			if rawPageNumber := r.URL.Query().Get("page"); len(rawPageNumber) != 0 {
				var err error
				if pageNumber, err = strconv.Atoi(rawPageNumber); err != nil {
					return &ResponseError{
						StatusCode: http.StatusBadRequest,
						Message:    "page invalid value",
					}
				}
			}

			ctx := context.WithValue(r.Context(), "paginator", pg.Copy(r.URL, int64(pageNumber), int64(limitCount), order))
			next.ServeHTTP(rw, r.WithContext(ctx))
			return nil
		})
	}
}
