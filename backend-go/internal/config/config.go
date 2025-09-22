package config

import (
	"os"
)

type Config struct {
	Port            string
	JWTSecret       string
	AllowedOrigins  string
	DatabaseURL     string
	PageSizeDefault int
}

func Load() Config {
	cfg := Config{
		Port:            getenv("PORT", "8080"),
		JWTSecret:       getenv("JWT_SECRET", "dev-secret"),
		AllowedOrigins:  getenv("ALLOWED_ORIGINS", "*"),
		DatabaseURL:     getenv("DATABASE_URL", ""),
		PageSizeDefault: 20,
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
