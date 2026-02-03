package main

import (
	"log"
	"os"
	"time"

	"analytics/internal/api/handlers"
	"analytics/internal/api/routes"
	"analytics/internal/db"
	"analytics/internal/repository"
	"analytics/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("stack.env")
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v\n", err)
	}

	databaseService := &db.DatabaseService{}
	pool := databaseService.GetPool()

	transactionRepo := repository.NewTransactionRepository(pool)
	categoryRepo := repository.NewCategoryRepository(pool)

	transactionAnalysisService := service.NewTransactionAnalysisService(
		transactionRepo,
		categoryRepo,
	)

	typeService := service.NewTypeService(transactionRepo)
	categoryService := service.NewCategoryService(categoryRepo, transactionRepo)

	transactionHandler := handlers.NewTransactionHandler(transactionRepo, transactionAnalysisService)
	typeHandler := handlers.NewTypeHandler(typeService)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo, categoryService)

	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins for development. Change for production!
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	router.SetTrustedProxies([]string{"172.16.0.0/12", "192.168.0.0/16"})

	routes.SetupRoutes(router, transactionHandler, typeHandler, categoryHandler)

	router.Run("0.0.0.0:1234")
}
