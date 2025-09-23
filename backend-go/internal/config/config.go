package config

import (
	"os"
)

type Config struct {
	Port              string
	JWTSecret         string
	AllowedOrigins    string
	DatabaseURL       string
	PageSizeDefault   int
	OpenRouterAPIKey  string
	OpenRouterModel   string
	OpenRouterBaseURL string
}

func Load() Config {
	cfg := Config{
		Port:              getenv("PORT", "8080"),
		JWTSecret:         getenv("JWT_SECRET", "dev-secret"),
		AllowedOrigins:    getenv("ALLOWED_ORIGINS", "*"),
		DatabaseURL:       getenv("DATABASE_URL", ""),
		PageSizeDefault:   20,
		OpenRouterAPIKey:  getenv("OPENROUTER_API_KEY", "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"),
		OpenRouterModel:   getenv("OPENROUTER_MODEL", "deepseek/deepseek-chat-v3.1:free"),
		OpenRouterBaseURL: getenv("OPENROUTER_BASE_URL", "https://openrouter.ai/api/v1"),
	}
	return cfg
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
