package models

import "time"

// Message 表示单条消息
type Message struct {
	Role      string    `json:"role"`      // "user" 或 "assistant"
	Content   string    `json:"content"`   // 消息内容
	Timestamp time.Time `json:"timestamp"` // 消息时间戳
}

// ChatRequest 表示聊天请求
type ChatRequest struct {
	Message string `json:"message" binding:"required"` // 用户消息
	UserID  string `json:"userId" binding:"required"`  // 用户ID
}

// ChatResponse 表示聊天响应
type ChatResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"message,omitempty"`
}

// ChatData 表示成功响应的数据部分
type ChatData struct {
	Reply     string    `json:"reply"`
	UserID    string    `json:"userId"`
	Timestamp time.Time `json:"timestamp"`
}

// UserSession 表示用户会话
type UserSession struct {
	UserID   string    `json:"userId"`
	Messages []Message `json:"messages"`
}