package services

import (
	"context"
	"errors"
	"fmt"
	"timekettle/chat/models"

	"github.com/sashabaranov/go-openai"
)

type OpenAIService struct {
	client *openai.Client
}

func NewOpenAIService(apiKey string) *OpenAIService {
	client := openai.NewClient(apiKey)
	return &OpenAIService{
		client: client,
	}
}

// GetChatResponse 获取 OpenAI 的聊天响应
func (s *OpenAIService) GetChatResponse(messages []models.Message) (string, error) {
	// 将消息转换为 OpenAI 格式
	var openaiMessages []openai.ChatCompletionMessage
	
	for _, msg := range messages {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 创建聊天完成请求
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo, // 改回 GPT-3.5-turbo，更稳定且便宜
		Messages:    openaiMessages,
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	// 调用 OpenAI API
	resp, err := s.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		// 添加更详细的错误信息
		return "", fmt.Errorf("OpenAI API 调用失败: %v", err)
	}

	// 检查响应
	if len(resp.Choices) == 0 {
		return "", errors.New("OpenAI 没有返回任何响应")
	}

	return resp.Choices[0].Message.Content, nil
}