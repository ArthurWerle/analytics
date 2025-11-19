package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseService struct {
	pool *pgxpool.Pool
}

func (db *DatabaseService) GetPool() *pgxpool.Pool {
	if db.pool != nil {
		return db.pool
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	db.pool = pool
	return pool
}

func (db *DatabaseService) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}
