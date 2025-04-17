package main

import (
	"context"
	"log"
	"os"

	"analytics/internal/api/handlers"
	"analytics/internal/api/routes"
	"analytics/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("stack.env")
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v\n", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("DATABASE_URL is not set")
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	transactionRepo := repository.NewTransactionRepository(conn)

	transactionHandler := handlers.NewTransactionHandler(transactionRepo)

	router := gin.Default()
	routes.SetupRoutes(router, transactionHandler)

	router.Run("0.0.0.0:1234")
}
