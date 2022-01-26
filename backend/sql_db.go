package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

func configureDatabasePath(env string) string {
	switch env {
	case "development", "production", "staging":
		return fmt.Sprintf("./_data/%s_%s.db", databasePrefix, env)
	default:
		log.Fatalf("invalid environment '%s'\n", env)
		return ""
	}
}

func runMigration(db *sqlx.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "./migrations",
	}

	n, err := migrate.Exec(db.DB, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return err
	}

	log.Printf("Applied %d migrations\n", n)
	return nil
}

func initializeDb() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalln(err)
	}

	if err := runMigration(db); err != nil {
		log.Fatalln(err)
	}

	return db
}

func wrapSqlResult(res sql.Result, customErrorMessage ...string) error {
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		errMessage := "Unable to process your submission. Please try again."
		if len(customErrorMessage) != 0 {
			errMessage = customErrorMessage[0]
		}
		return &ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    errMessage,
		}
	}
	return nil
}

type FilterFunc func(*http.Request, string, *sq.SelectBuilder) error

func customSelectFilters(filters map[string]FilterFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			selectQuery := sq.Select()
			query := r.URL.Query()
			for targetQueryName, queryFunc := range filters {
				queryVal := query.Get(targetQueryName)
				if err := queryFunc(r, queryVal, &selectQuery); err != nil {
					return err
				}
			}
			ctx := context.WithValue(r.Context(), "selectQuery", selectQuery)
			next.ServeHTTP(rw, r.WithContext(ctx))
			return nil
		})
	}
}

type Predicate interface {
	ToSql() (string, []interface{}, error)
}
