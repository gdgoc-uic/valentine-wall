package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
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

	fmt.Printf("Applied %d migrations\n", n)
	return nil
}

func initializeDb() *sqlx.DB {
	_, fileDbErr := os.Stat(databasePath)
	db, err := sqlx.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalln(err)
	}

	if errors.Is(fileDbErr, os.ErrNotExist) {
		if err := runMigration(db); err != nil {
			log.Fatalln(err)
		}
	}

	return db
}
