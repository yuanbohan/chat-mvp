package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuanbohan/chat-mvp/config"
	"github.com/yuanbohan/chat-mvp/handlers"
	"github.com/yuanbohan/chat-mvp/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	if cfg.OpenAIAPIKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Initialize services
	openaiService := services.NewOpenAIService(cfg.OpenAIAPIKey)

	// Initialize handlers
	chatHandler := handlers.NewChatHandler(openaiService)

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Serve static files
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")

	// API routes
	api := r.Group("/api")
	{
		api.POST("/chat", chatHandler.Chat)
		api.GET("/history", chatHandler.GetHistory)
		api.DELETE("/history", chatHandler.ClearHistory)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "chat-mvp",
		})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}