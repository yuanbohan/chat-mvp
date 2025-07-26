package config

import (
	"os"
)

// Config 应用配置结构
type Config struct {
	OpenAIAPIKey string
	Port         string
	Environment  string
}

// Load 加载配置
func Load() *Config {
	config := &Config{
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
		Port:         getEnv("PORT", "8080"),
		Environment:  getEnv("ENVIRONMENT", "development"),
	}

	if config.OpenAIAPIKey == "" {
		panic("OPENAI_API_KEY environment variable is required")
	}

	return config
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
