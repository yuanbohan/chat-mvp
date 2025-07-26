package models

import "time"

// Message represents a chat message
type Message struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	Role      string    `json:"role"` // "user" or "assistant"
	Timestamp time.Time `json:"timestamp"`
}

// ChatRequest represents the incoming chat request
type ChatRequest struct {
	Message string `json:"message" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
}

// ChatResponse represents the response from the chat API
type ChatResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// Session represents a user's chat session
type Session struct {
	UserID   string    `json:"user_id"`
	Messages []Message `json:"messages"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}