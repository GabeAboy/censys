package server

import (
	"fmt"
	"log"

	"censys/internal/database"
	"censys/internal/handlers"
	"censys/internal/services"
	"censys/pkg/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "censys/docs"
)

func Run() error {
	// Load config
	cfg := config.Load()

	// Initialize database
	db, err := database.NewDB(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// TODO Use golang-migrate
	if err := db.InitSchema(); err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	// Initialize repositories
	assetRepo := database.NewAssetRepository(db)
	portRepo := database.NewPortRepository(db)
	tagRepo := database.NewTagRepository(db)

	// Initialize services
	assetService := services.NewAssetService(db, assetRepo, portRepo, tagRepo)

	// Initialize handlers
	assetHandler := handlers.NewAssetHandler(assetRepo, assetService)

	// Setup router
	router := setupRouter(assetHandler)

	// Start server
	log.Printf("Starting Censys API server on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func setupRouter(assetHandler *handlers.AssetHandler) *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(corsMiddleware())

	// Health check
	router.GET("/health", healthCheck)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	v1 := router.Group("/api/v1")
	{
		assets := v1.Group("/assets")
		{
			assets.GET("", assetHandler.GetAssetList)
			assets.GET("/count", assetHandler.GetAssetCount)
			assets.POST("", assetHandler.CreateAsset)
			assets.DELETE("/:id", assetHandler.DeleteAsset)
			assets.POST("/:id/tags", assetHandler.CreateAssetTag)
		}
	}

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "healthy"})
}
