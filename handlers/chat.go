package handlers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yuanbohan/chat-mvp/models"
	"github.com/yuanbohan/chat-mvp/services"
)

// ChatHandler handles chat-related requests
type ChatHandler struct {
	openaiService *services.OpenAIService
	sessions      map[string]*models.Session
	sessionsMutex sync.RWMutex
}

// NewChatHandler creates a new chat handler
func NewChatHandler(openaiService *services.OpenAIService) *ChatHandler {
	return &ChatHandler{
		openaiService: openaiService,
		sessions:      make(map[string]*models.Session),
	}
}

// Chat handles the chat API endpoint
func (h *ChatHandler) Chat(c *gin.Context) {
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	// Get or create session
	_ = h.getOrCreateSession(req.UserID)

	// Create user message
	userMessage := models.Message{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Content:   req.Message,
		Role:      "user",
		Timestamp: time.Now(),
	}

	// Add user message to session
	h.addMessageToSession(req.UserID, userMessage)

	// Get all messages for context
	messages := h.getSessionMessages(req.UserID)

	// Get response from OpenAI
	response, err := h.openaiService.Chat(c.Request.Context(), messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to get AI response: %v", err),
		})
		return
	}

	// Create assistant message
	assistantMessage := models.Message{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Content:   response,
		Role:      "assistant",
		Timestamp: time.Now(),
	}

	// Add assistant message to session
	h.addMessageToSession(req.UserID, assistantMessage)

	// Return response
	c.JSON(http.StatusOK, models.ChatResponse{
		Success:   true,
		Message:   response,
		Timestamp: time.Now(),
	})
}

// GetHistory returns the chat history for a user
func (h *ChatHandler) GetHistory(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "user_id is required",
		})
		return
	}

	messages := h.getSessionMessages(userID)
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"messages": messages,
	})
}

// ClearHistory clears the chat history for a user
func (h *ChatHandler) ClearHistory(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "user_id is required",
		})
		return
	}

	h.sessionsMutex.Lock()
	delete(h.sessions, userID)
	h.sessionsMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "History cleared successfully",
	})
}

// getOrCreateSession gets or creates a session for a user
func (h *ChatHandler) getOrCreateSession(userID string) *models.Session {
	h.sessionsMutex.Lock()
	defer h.sessionsMutex.Unlock()

	session, exists := h.sessions[userID]
	if !exists {
		session = &models.Session{
			UserID:   userID,
			Messages: make([]models.Message, 0),
			Created:  time.Now(),
			Updated:  time.Now(),
		}
		h.sessions[userID] = session
	}

	return session
}

// addMessageToSession adds a message to a user's session
func (h *ChatHandler) addMessageToSession(userID string, message models.Message) {
	h.sessionsMutex.Lock()
	defer h.sessionsMutex.Unlock()

	if session, exists := h.sessions[userID]; exists {
		session.Messages = append(session.Messages, message)
		session.Updated = time.Now()
	}
}

// getSessionMessages returns all messages for a user's session
func (h *ChatHandler) getSessionMessages(userID string) []models.Message {
	h.sessionsMutex.RLock()
	defer h.sessionsMutex.RUnlock()

	if session, exists := h.sessions[userID]; exists {
		return session.Messages
	}
	return []models.Message{}
}