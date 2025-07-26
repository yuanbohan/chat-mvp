package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAIAPIKey string
	Port         string
}

func LoadConfig() *Config {
	// 尝试加载 .env 文件，如果不存在也不会报错
	err := godotenv.Load()
	if err != nil {
		log.Println("没有找到 .env 文件，使用环境变量")
	}

	config := &Config{
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
		Port:         getEnv("PORT", "8080"),
	}

	if config.OpenAIAPIKey == "" {
		log.Fatal("OPENAI_API_KEY 环境变量未设置")
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}