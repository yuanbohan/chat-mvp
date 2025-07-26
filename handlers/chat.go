package handlers

import (
	"context"
	"net/http"
	"time"

	"chat-mvp/models"
	"chat-mvp/services"

	"github.com/gin-gonic/gin"
)

// ChatHandler 聊天处理器结构体
type ChatHandler struct {
	openaiService *services.OpenAIService
}

// NewChatHandler 创建新的聊天处理器
func NewChatHandler(openaiService *services.OpenAIService) *ChatHandler {
	return &ChatHandler{
		openaiService: openaiService,
	}
}

// Chat 处理聊天请求
func (h *ChatHandler) Chat(c *gin.Context) {
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "请求格式错误: " + err.Error(),
		})
		return
	}

	// 验证必要字段
	if req.Message == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "消息内容不能为空",
		})
		return
	}

	if req.UserID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "用户ID不能为空",
		})
		return
	}

	// 创建超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 调用 OpenAI 服务
	response, err := h.openaiService.Chat(ctx, req.UserID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "AI 服务暂时不可用，请稍后再试",
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, models.ChatResponse{
		Success:   true,
		Message:   response,
		Timestamp: time.Now(),
		UserID:    req.UserID,
	})
}

// GetHistory 获取用户聊天历史
func (h *ChatHandler) GetHistory(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "用户ID不能为空",
		})
		return
	}

	history := h.openaiService.GetUserHistory(userID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"history": history,
	})
}

// ClearHistory 清除用户聊天历史
func (h *ChatHandler) ClearHistory(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "请求格式错误: " + err.Error(),
		})
		return
	}

	h.openaiService.ClearUserHistory(req.UserID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "聊天历史已清除",
	})
}

// Health 健康检查接口
func (h *ChatHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now(),
	})
}
