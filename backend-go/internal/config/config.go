package config

import (
    "os"
)

type Config struct {
    Port             string
    JWTSecret        string
    AllowedOrigins   string
}

func Load() Config {
    cfg := Config{
        Port:           getenv("PORT", "8080"),
        JWTSecret:      getenv("JWT_SECRET", "dev-secret"),
        AllowedOrigins: getenv("ALLOWED_ORIGINS", "*"),
    }
    return cfg
}

func getenv(key, def string) string {
    v := os.Getenv(key)
    if v == "" { return def }
    return v
}


