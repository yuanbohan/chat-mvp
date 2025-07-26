package main

import (
	"log"

	"chat-mvp/config"
	"chat-mvp/handlers"
	"chat-mvp/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化 OpenAI 服务
	openaiService := services.NewOpenAIService(cfg.OpenAIAPIKey)

	// 初始化处理器
	chatHandler := handlers.NewChatHandler(openaiService)

	// 初始化 Gin 路由器
	router := gin.New()

	// 添加中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 添加 CORS 中间件
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 静态文件服务
	router.Static("/static", "./static")
	router.LoadHTMLGlob("static/*.html")

	// 主页路由
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Chat MVP - OpenAI 聊天应用",
		})
	})

	// API 路由组
	api := router.Group("/api")
	{
		api.POST("/chat", chatHandler.Chat)
		api.GET("/history", chatHandler.GetHistory)
		api.POST("/clear-history", chatHandler.ClearHistory)
		api.GET("/health", chatHandler.Health)
	}

	// 启动服务器
	port := ":" + cfg.Port
	log.Printf("服务器启动在端口 %s", port)
	log.Printf("访问地址: http://127.0.0.1%s", port)

	if err := router.Run(port); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}
