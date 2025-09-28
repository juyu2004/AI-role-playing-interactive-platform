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
	BaiduAppID        string
	BaiduAPIKey       string
	BaiduSecretKey    string
}

func Load() Config {
	cfg := Config{
		Port:           getenv("PORT", "8080"),
		JWTSecret:      getenv("JWT_SECRET", "dev-secret"),
		AllowedOrigins: getenv("ALLOWED_ORIGINS", "*"),
		// 默认本地 PostgreSQL 连接串。建议通过环境变量 DATABASE_URL 覆盖。
		// 示例：postgres://用户名:密码@主机:端口/数据库名?sslmode=disable
		// 本地示例：postgres://postgres:123456@localhost:5432/mydb?sslmode=disable
		DatabaseURL:       getenv("DATABASE_URL", "postgres://postgres:123456@localhost:5432/mydb?sslmode=disable"),
		PageSizeDefault:   20,
		OpenRouterAPIKey:  getenv("OPENROUTER_API_KEY", "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"),
		OpenRouterModel:   getenv("OPENROUTER_MODEL", "nvidia/nemotron-nano-9b-v2:free"),
		OpenRouterBaseURL: getenv("OPENROUTER_BASE_URL", "https://openrouter.ai/api/v1"),
		// 默认内置百度智能云语音服务凭据（可直接使用）。
		BaiduAppID:     getenv("BAIDU_APP_ID", "7080797"),
		BaiduAPIKey:    getenv("BAIDU_API_KEY", "By5O9RHtXiycynI7OUMrvgVI"),
		BaiduSecretKey: getenv("BAIDU_SECRET_KEY", "UU9jUJgxR5XDe1iuuhwPbSFoZa8Zm70g"),
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
