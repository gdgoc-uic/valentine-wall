package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpgx"
	migrate "github.com/rubenv/sql-migrate"
)

type DatabaseDriver interface {
	Init(string) (*sqlx.DB, error)
}

type PostgresDB struct{}

func (sl *PostgresDB) Init(conn string) (*sqlx.DB, error) {
	log.Printf("connecting to %s...\n", conn)
	dbDriver := "pgx"
	if newrelicApp != nil {
		dbDriver = "nrpgx"
	}
	return sqlx.Open(dbDriver, conn)
}

type SQLiteDB struct{}

func (sl *SQLiteDB) Init(dbPath string) (*sqlx.DB, error) {
	if _, err := os.Stat(dataDirPath); errors.Is(err, os.ErrNotExist) {
		log.Printf("data directory not found. creating one...")
		if err := os.Mkdir(dataDirPath, 0777); err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		log.Printf("database not found. creating one...")
		if _, err := os.OpenFile(dbPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777); err != nil {
			log.Fatalln(err)
		}
	}

	log.Printf("connecting to %s...\n", dbPath)
	return sqlx.Open("sqlite3", dbPath)
}

var dbDrivers = map[string]DatabaseDriver{
	"sqlite3":  &SQLiteDB{},
	"postgres": &PostgresDB{},
}

func configureDatabasePath(env string) (string, string) {
	switch env {
	case "development", "production", "staging":
		username := ""
		password := ""

		if gotPostgresUsername, exists := os.LookupEnv("POSTGRES_USER"); exists {
			username = "user=" + gotPostgresUsername
		}

		if gotPostgresPassword, exists := os.LookupEnv("POSTGRES_PASSWORD"); exists {
			password = "password=" + gotPostgresPassword
		}

		return "postgres", fmt.Sprintf("dbname=%s_%s %s %s sslmode=disable", databasePrefix, env, username, password)
	default:
		log.Fatalf("invalid environment '%s'\n", env)
		return "", ""
	}
}

func runMigration(db *sqlx.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "./migrations",
	}

	n, err := migrate.Exec(db.DB, databaseDriver, migrations, migrate.Up)
	if err != nil {
		return err
	}

	log.Printf("Applied %d migrations\n", n)
	return nil
}

func initializeDb() *sqlx.DB {
	db, err := dbDrivers[databaseDriver].Init(databasePath)
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

type selectQueryKey struct{}

func getSelectQueryFromReq(r *http.Request) *sq.SelectBuilder {
	sQuery, ok := r.Context().Value(selectQueryKey{}).(*sq.SelectBuilder)
	if !ok {
		return nil
	}
	return sQuery
}

func injectSelectQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(selectQueryKey{}).(*sq.SelectBuilder)
		if !ok {
			selectQuery := psql.Select()
			ctx := context.WithValue(r.Context(), selectQueryKey{}, &selectQuery)
			next.ServeHTTP(rw, r.WithContext(ctx))
			return
		}
		next.ServeHTTP(rw, r)
	})
}
