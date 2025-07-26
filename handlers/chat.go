package handlers

import (
	"log"
	"net/http"
	"sync"
	"time"
	"timekettle/chat/models"
	"timekettle/chat/services"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	openaiService *services.OpenAIService
	userSessions  map[string]*models.UserSession
	mutex         sync.RWMutex
}

func NewChatHandler(openaiService *services.OpenAIService) *ChatHandler {
	return &ChatHandler{
		openaiService: openaiService,
		userSessions:  make(map[string]*models.UserSession),
		mutex:         sync.RWMutex{},
	}
}

// HandleChat 处理聊天请求
func (h *ChatHandler) HandleChat(c *gin.Context) {
	var req models.ChatRequest
	
	log.Printf("收到聊天请求")
	
	// 绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("请求格式错误: %v", err)
		c.JSON(http.StatusBadRequest, models.ChatResponse{
			Status: "error",
			Error:  "请求格式错误: " + err.Error(),
		})
		return
	}

	log.Printf("用户 %s 发送消息: %s", req.UserID, req.Message)

	// 获取或创建用户会话
	h.mutex.Lock()
	session, exists := h.userSessions[req.UserID]
	if !exists {
		session = &models.UserSession{
			UserID:   req.UserID,
			Messages: []models.Message{},
		}
		h.userSessions[req.UserID] = session
		log.Printf("为用户 %s 创建新会话", req.UserID)
	}
	h.mutex.Unlock()

	// 添加用户消息到会话
	userMessage := models.Message{
		Role:      "user",
		Content:   req.Message,
		Timestamp: time.Now(),
	}
	
	h.mutex.Lock()
	session.Messages = append(session.Messages, userMessage)
	h.mutex.Unlock()

	log.Printf("开始调用 OpenAI API...")
	// 调用 OpenAI API
	reply, err := h.openaiService.GetChatResponse(session.Messages)
	if err != nil {
		log.Printf("OpenAI API 调用失败: %v", err)
		c.JSON(http.StatusInternalServerError, models.ChatResponse{
			Status: "error",
			Error:  "获取 AI 响应失败: " + err.Error(),
		})
		return
	}

	log.Printf("OpenAI API 调用成功，回复: %s", reply)

	// 添加 AI 响应到会话
	assistantMessage := models.Message{
		Role:      "assistant",
		Content:   reply,
		Timestamp: time.Now(),
	}
	
	h.mutex.Lock()
	session.Messages = append(session.Messages, assistantMessage)
	h.mutex.Unlock()

	// 返回响应
	c.JSON(http.StatusOK, models.ChatResponse{
		Status: "success",
		Data: models.ChatData{
			Reply:     reply,
			UserID:    req.UserID,
			Timestamp: time.Now(),
		},
	})
	
	log.Printf("成功处理用户 %s 的聊天请求", req.UserID)
}

// GetUserHistory 获取用户聊天历史（可选接口）
func (h *ChatHandler) GetUserHistory(c *gin.Context) {
	userID := c.Query("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, models.ChatResponse{
			Status: "error",
			Error:  "缺少 userId 参数",
		})
		return
	}

	h.mutex.RLock()
	session, exists := h.userSessions[userID]
	h.mutex.RUnlock()

	if !exists {
		c.JSON(http.StatusOK, models.ChatResponse{
			Status: "success",
			Data:   []models.Message{},
		})
		return
	}

	c.JSON(http.StatusOK, models.ChatResponse{
		Status: "success",
		Data:   session.Messages,
	})
}