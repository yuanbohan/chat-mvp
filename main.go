package main

import (
	"log"
	"timekettle/chat/config"
	"timekettle/chat/handlers"
	"timekettle/chat/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 创建 OpenAI 服务
	openaiService := services.NewOpenAIService(cfg.OpenAIAPIKey)

	// 创建聊天处理器
	chatHandler := handlers.NewChatHandler(openaiService)

	// 创建 Gin 路由器
	router := gin.Default()

	// 设置静态文件服务
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")

	// API 路由
	api := router.Group("/api")
	{
		api.POST("/chat", chatHandler.HandleChat)
		api.GET("/history", chatHandler.GetUserHistory)
	}

	// 兼容原接口路径
	router.POST("/chat", chatHandler.HandleChat)

	// 启动服务器
	log.Printf("服务器启动在端口 %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}