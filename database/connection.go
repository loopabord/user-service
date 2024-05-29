package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Import the pq driver
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var db *bun.DB

func Connect() error {
	databaseUser := os.Getenv("DATABASE_USER")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseURL := os.Getenv("DATABASE_URL")
	database := os.Getenv("DATABASE")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", databaseUser, databasePassword, databaseURL, database)
	// dsn := "unix://user:pass@dbname/var/run/postgresql/.s.PGSQL.5432"
	pgconn := pgdriver.NewConnector(pgdriver.WithDSN(dsn))
	sqldb := sql.OpenDB(pgconn)

	db = bun.NewDB(sqldb, pgdialect.New())

	// Enable query debugging
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return nil
}

func Initialize() error {
	databaseUser := os.Getenv("DATABASE_USER")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseURL := os.Getenv("DATABASE_URL")
	database := os.Getenv("DATABASE")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", databaseUser, databasePassword, databaseURL, database)
	log.Println(dsn)
	sqldb, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}
	defer sqldb.Close()

	// Check if the user database exists
	var exists bool
	err = sqldb.QueryRow("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = 'artist')").Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if user database exists: %w", err)
	}

	// Create the user database if it doesn't exist
	if !exists {
		_, err = sqldb.Exec(`CREATE DATABASE "artist"`)
		if err != nil {
			return fmt.Errorf("failed to create user database: %w", err)
		}
		log.Println("Database artist created successfully.")
	}

	return nil
}

// Ensure you close the db connection when the application shuts down
func Close() error {
	return db.Close()
}
