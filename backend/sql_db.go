package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

func configureDatabasePath(env string) string {
	switch env {
	case "development", "production", "staging":
		return filepath.Join(dataDirPath, fmt.Sprintf("%s_%s.db", databasePrefix, env))
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
	if _, err := os.Stat(dataDirPath); errors.Is(err, os.ErrNotExist) {
		log.Printf("data directory not found. creating one...")
		if err := os.Mkdir(dataDirPath, 0777); err != nil {
			log.Fatalln(err)
		}
	}

	log.Printf("connecting to %s...\n", databasePath)
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

type Predicate interface {
	ToSql() (string, []interface{}, error)
}

func injectSelectQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value("selectQuery").(*sq.SelectBuilder)
		if !ok {
			selectQuery := sq.Select()
			ctx := context.WithValue(r.Context(), "selectQuery", &selectQuery)
			next.ServeHTTP(rw, r.WithContext(ctx))
			return
		}
		next.ServeHTTP(rw, r)
	})
}
