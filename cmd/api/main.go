package main

import (
	"context"
	"log"
	"os"

	"analytics/internal/api/handlers"
	"analytics/internal/api/routes"
	"analytics/internal/repository"
	"analytics/internal/service"

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
	categoryRepo := repository.NewCategoryRepository(conn)
	typeRepo := repository.NewTypeRepository(conn)
	recurringRepo := repository.NewRecurringTransactionRepository(conn)

	transactionAnalysisService := service.NewTransactionAnalysisService(
		transactionRepo,
		recurringRepo,
		categoryRepo,
		typeRepo,
	)

	transactionHandler := handlers.NewTransactionHandler(transactionRepo, transactionAnalysisService)

	// Set Gin to release mode in production
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.SetTrustedProxies([]string{"172.16.0.0/12", "192.168.0.0/16"}) // Docker network ranges

	routes.SetupRoutes(router, transactionHandler)

	router.Run("0.0.0.0:1234")
}
