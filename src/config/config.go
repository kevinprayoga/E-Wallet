package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	"log"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("POSTGRES_DB")

	log.Printf("DB_HOST=%s, DB_USER=%s, DB_NAME=%s", host, user, dbname)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Second * 1)

	if err := HealthCheck(db); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func HealthCheck(db *sql.DB) error {
	return db.Ping()
}
