package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type DatabaseService struct {
	conn *pgx.Conn
}

func (db *DatabaseService) GetConnection() *pgx.Conn {
	if db.conn != nil {
		return db.conn
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("DATABASE_URL is not set")
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	db.conn = conn
	return conn
}

func (db *DatabaseService) Close() error {
	if db.conn != nil {
		return db.conn.Close(context.Background())
	}
	return nil
}
