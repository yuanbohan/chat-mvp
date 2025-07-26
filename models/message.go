package models

import "time"

// Message 消息结构体
type Message struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	Role      string    `json:"role"` // "user" 或 "assistant"
	Timestamp time.Time `json:"timestamp"`
}

// ChatRequest 聊天请求结构体
type ChatRequest struct {
	Message string `json:"message" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
}

// ChatResponse 聊天响应结构体
type ChatResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// ErrorResponse 错误响应结构体
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
