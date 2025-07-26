package services

import (
	"context"
	"fmt"
	"sync"

	"chat-mvp/models"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAIService OpenAI 服务结构体
type OpenAIService struct {
	client   *openai.Client
	sessions map[string][]models.Message // 用户会话历史
	mutex    sync.RWMutex
}

// NewOpenAIService 创建新的 OpenAI 服务实例
func NewOpenAIService(apiKey string) *OpenAIService {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &OpenAIService{
		client:   &client,
		sessions: make(map[string][]models.Message),
		mutex:    sync.RWMutex{},
	}
}

// Chat 处理聊天请求
func (s *OpenAIService) Chat(ctx context.Context, userID, message string) (string, error) {
	// 获取用户会话历史
	s.mutex.Lock()
	history := s.sessions[userID]

	// 添加用户消息到历史
	userMessage := models.Message{
		UserID:  userID,
		Content: message,
		Role:    "user",
	}
	history = append(history, userMessage)
	s.sessions[userID] = history
	s.mutex.Unlock()

	// 构建 OpenAI 消息格式
	messages := s.buildMessages(history)

	// 调用 OpenAI API
	chatCompletion, err := s.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    openai.ChatModelGPT4o,
	})

	if err != nil {
		return "", fmt.Errorf("OpenAI API 调用失败: %w", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return "", fmt.Errorf("OpenAI API 返回空响应")
	}

	// 获取助手回复
	assistantReply := chatCompletion.Choices[0].Message.Content

	// 保存助手回复到历史
	s.mutex.Lock()
	assistantMessage := models.Message{
		UserID:  userID,
		Content: assistantReply,
		Role:    "assistant",
	}
	s.sessions[userID] = append(s.sessions[userID], assistantMessage)
	s.mutex.Unlock()

	return assistantReply, nil
}

// buildMessages 构建 OpenAI API 消息格式
func (s *OpenAIService) buildMessages(history []models.Message) []openai.ChatCompletionMessageParamUnion {
	// 添加系统消息
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("你是一个友好且有帮助的 AI 助手，基于 GPT-4o 模型。请用中文回复用户的问题。"),
	}

	// 限制历史消息数量，避免超过 token 限制
	maxHistory := 20
	start := 0
	if len(history) > maxHistory {
		start = len(history) - maxHistory
	}

	// 添加历史消息
	for _, msg := range history[start:] {
		if msg.Role == "user" {
			messages = append(messages, openai.UserMessage(msg.Content))
		} else if msg.Role == "assistant" {
			messages = append(messages, openai.AssistantMessage(msg.Content))
		}
	}

	return messages
}

// GetUserHistory 获取用户聊天历史
func (s *OpenAIService) GetUserHistory(userID string) []models.Message {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if history, exists := s.sessions[userID]; exists {
		// 返回副本，避免并发修改
		result := make([]models.Message, len(history))
		copy(result, history)
		return result
	}
	return []models.Message{}
}

// ClearUserHistory 清除用户聊天历史
func (s *OpenAIService) ClearUserHistory(userID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.sessions, userID)
}
