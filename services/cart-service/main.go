package main

import (
	"cart-service/handlers"
	"cart-service/middleware"
	"cart-service/utils"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Initialize Redis
	if err := utils.InitRedis(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Create Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "cart-service",
		})
	})

	// API routes (protected)
	api := router.Group("/api/cart")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("", handlers.GetCart)
		api.POST("/items", handlers.AddItem)
		api.PUT("/items/:product_id", handlers.UpdateItem)
		api.DELETE("/items/:product_id", handlers.RemoveItem)
		api.DELETE("", handlers.ClearCart)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5003"
	}

	fmt.Printf("🚀 Cart Service running on port %s\n", port)
	router.Run(fmt.Sprintf(":%s", port))
}