package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/blevesearch/bleve"
	"github.com/jmoiron/sqlx"
)

var paginatorCtxKey = struct{}{}

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
	Len        int            `json:"-"`
}

type Paginator struct {
	Path      *url.URL
	Page      int64
	Limit     int64
	OrderKey  string
	Order     string
	TableName string
}

func (pg *Paginator) Copy(path *url.URL, page int64, limit int64, order string) *Paginator {
	return &Paginator{
		Path:     path,
		OrderKey: pg.OrderKey,
		Page:     page,
		Limit:    limit,
		Order:    order,
	}
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

func (pg *Paginator) Load(source PaginatorSource) (*PaginatedResponse, error) {
	count, err := source.Count()
	if err != nil {
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

	// fmt.Printf("{order: %s, limit: %d, page: %d, OrderKey: %s}\n", pg.Order, pg.Limit, pg.Page, pg.OrderKey)
	results, resultsLen, err := source.Fetch(pg.Page, pg.Limit, pg.OrderKey, pg.Order)
	if err != nil {
		return nil, err
	}

	orderQuery := pg.Order
	if len(pg.OrderKey) != 0 {
		orderQuery = pg.OrderKey + "," + orderQuery
	}

	resp := &PaginatedResponse{
		Links: PaginatedLinks{
			First: *generatePaginateUrl(pg.Path, 1, pg.Limit, orderQuery),
			Last:  *generatePaginateUrl(pg.Path, lastPageNumber, pg.Limit, orderQuery),
		},
		Page:       currentPageNumber,
		TotalCount: count,
		PerPage:    pg.Limit,
		PageCount:  lastPageNumber,
		Data:       results,
		Len:        resultsLen,
	}

	if currentPageNumber+1 <= lastPageNumber {
		resp.Links.Next = generatePaginateUrl(pg.Path, currentPageNumber+1, pg.Limit, orderQuery)
	}

	if currentPageNumber > 1 {
		resp.Links.Previous = generatePaginateUrl(pg.Path, currentPageNumber-1, pg.Limit, orderQuery)
	}

	return resp, nil
}

func pagination(pg *Paginator) func(http.Handler) http.Handler {
	filterMiddleware := customFilters(map[string]FilterFunc{
		"order": func(r *http.Request, ctx context.Context, f Filter) error {
			if !f.Exists || len(f.Value) == 0 {
				return nil
			}

			pg := ctx.Value(paginatorCtxKey).(*Paginator)
			splitted := strings.Split(f.Value, ",")
			lastIndex := len(splitted) - 1
			if len(splitted) > 2 {
				return &ResponseError{
					StatusCode: http.StatusUnprocessableEntity,
					Message:    "ordering by specific column must be column name follow by a single comma and order (asc/desc)",
				}
			} else if splitted[lastIndex] != "desc" && splitted[lastIndex] != "asc" {
				return &ResponseError{
					StatusCode: http.StatusBadRequest,
					Message:    "order invalid value",
				}
			} else if len(splitted) == 2 {
				pg.OrderKey = splitted[0]
				pg.Order = splitted[1]
			} else {
				pg.Order = f.Value
			}
			return nil
		},
		"limit": func(r *http.Request, ctx context.Context, f Filter) error {
			if !f.Exists || len(f.Value) == 0 {
				return nil
			}

			pg := ctx.Value(paginatorCtxKey).(*Paginator)
			limitCount, err := strconv.Atoi(f.Value)
			if err != nil {
				return &ResponseError{
					StatusCode: http.StatusBadRequest,
					Message:    "limit invalid value",
				}
			}

			pg.Limit = int64(limitCount)
			return nil
		},
		"page": func(r *http.Request, ctx context.Context, f Filter) error {
			if !f.Exists || len(f.Value) == 0 {
				return nil
			}

			pg := ctx.Value(paginatorCtxKey).(*Paginator)
			pageNumber, err := strconv.Atoi(f.Value)
			if err != nil {
				return &ResponseError{
					StatusCode: http.StatusBadRequest,
					Message:    "page invalid value",
				}
			}

			pg.Page = int64(pageNumber)
			return nil
		},
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), paginatorCtxKey, pg.Copy(r.URL, 1, 10, "asc"))
			filterMiddleware(next).ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

type PaginatorSource interface {
	Count() (int64, error)
	Fetch(page, limit int64, orderKey, order string) (interface{}, int, error)
}

type DatabasePaginatorSource struct {
	BaseQuery sq.SelectBuilder
	DataQuery sq.SelectBuilder
	DB        *sqlx.DB
	Converter func(*sqlx.Rows) (interface{}, error)
}

func (src *DatabasePaginatorSource) Count() (int64, error) {
	count := int64(0)
	if countQuery, args, err := src.BaseQuery.Column("count(*)").ToSql(); err != nil {
		return 0, err
	} else if rows, err := src.DB.Query(countQuery, args...); err != nil {
		return 0, err
	} else if rows.Next() {
		rows.Scan(&count)
		rows.Close()
	}
	return count, nil
}

func (src *DatabasePaginatorSource) Fetch(page, limit int64, orderKey, order string) (interface{}, int, error) {
	finalDataQuery, args, err := src.DataQuery.
		OrderBy(orderKey + " " + order).
		Limit(uint64(limit)).
		Offset(uint64(page-1) * uint64(limit)).ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := src.DB.Queryx(finalDataQuery, args...)
	if err != nil {
		return nil, 0, &ResponseError{StatusCode: http.StatusNotFound, WError: err}
	}
	defer rows.Close()

	results := []interface{}{}
	for rows.Next() {
		gotData, err := src.Converter(rows)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, gotData)
	}
	return results, len(results), nil
}

type BlevePaginatorSource struct {
	Index         bleve.Index
	SearchRequest *bleve.SearchRequest
}

func (src *BlevePaginatorSource) Count() (int64, error) {
	resp, err := src.Index.Search(src.SearchRequest)
	if err != nil {
		return 0, err
	}
	return int64(resp.Total), nil
}

func (src *BlevePaginatorSource) Fetch(page, limit int64, orderKey, order string) (interface{}, int, error) {
	finalOrder := orderKey
	if order == "desc" {
		finalOrder = "-" + finalOrder
	}

	src.SearchRequest.Size = int(limit)
	src.SearchRequest.From = int((page - 1) * limit)
	src.SearchRequest.Explain = false
	src.SearchRequest.SortBy([]string{finalOrder, "-_score", "_id"})
	resp, err := src.Index.Search(src.SearchRequest)
	if err != nil {
		return nil, 0, err
	}

	results := make([]interface{}, len(resp.Hits))
	for i, result := range resp.Hits {
		results[i] = result
	}
	return results, len(results), nil
}

type PipePaginatorSource struct {
	Source   PaginatorSource
	PipeFunc func([]interface{}) ([]interface{}, error)
}

func (src *PipePaginatorSource) Count() (int64, error) {
	return src.Source.Count()
}

func (src *PipePaginatorSource) Fetch(page, limit int64, orderKey, order string) (interface{}, int, error) {
	initialResults, _, err := src.Source.Fetch(page, limit, orderKey, order)
	if err != nil {
		return nil, 0, err
	}
	finalResults, err := src.PipeFunc(initialResults.([]interface{}))
	if err != nil {
		return nil, 0, err
	}
	return finalResults, len(finalResults), nil
}

type ArrayPaginatorSource struct {
	Arr interface{}
}

func (src *ArrayPaginatorSource) Count() (int64, error) {
	v := reflect.ValueOf(src.Arr)
	return int64(v.Len()), nil
}

func (src *ArrayPaginatorSource) Fetch(page, limit int64, orderKey, order string) (interface{}, int, error) {
	arr := reflect.ValueOf(src.Arr)
	if arr.Kind() != reflect.Slice {
		return nil, 0, fmt.Errorf("not a slice")
	}
	slice := reflect.MakeSlice(arr.Type(), 0, int(limit))
	j := 0
	for i := int(page-1) * int(limit); i < arr.Len(); i++ {
		if j >= int(limit) {
			break
		}
		slice = reflect.Append(slice, arr.Index(i))
		j++
	}
	return slice.Interface(), slice.Len(), nil
}

func getPaginatorFromReq(r *http.Request) *Paginator {
	return r.Context().Value(paginatorCtxKey).(*Paginator)
}
