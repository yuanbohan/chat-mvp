package services

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/yuanbohan/chat-mvp/models"
)

type OpenAIService struct {
	client openai.Client
}

// NewOpenAIService creates a new OpenAI service instance
func NewOpenAIService(apiKey string) *OpenAIService {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	
	return &OpenAIService{
		client: client,
	}
}

// Chat sends a message to OpenAI and returns the response
func (s *OpenAIService) Chat(ctx context.Context, messages []models.Message) (string, error) {
	// Convert our messages to OpenAI format
	openaiMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))
	
	for _, msg := range messages {
		switch msg.Role {
		case "user":
			openaiMessages = append(openaiMessages, openai.UserMessage(msg.Content))
		case "assistant":
			openaiMessages = append(openaiMessages, openai.AssistantMessage(msg.Content))
		}
	}

	// Create chat completion request
	response, err := s.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages:    openaiMessages,
		Model:       openai.ChatModelGPT4o,
		MaxTokens:   openai.Int(1000),
		Temperature: openai.Float(0.7),
	})
	
	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return response.Choices[0].Message.Content, nil
}