package main

import (
	"database/sql"
	"log"
	"net/http"

	"ai-role-playing-platform/backend-go/internal/app"
	"ai-role-playing-platform/backend-go/internal/config"
	"ai-role-playing-platform/backend-go/internal/handlers"
	"ai-role-playing-platform/backend-go/internal/middleware"
	rdb "ai-role-playing-platform/backend-go/internal/repo/db"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()
	var sqlDB *sql.DB
	if cfg.DatabaseURL != "" {
		if db, err := rdb.Open(cfg.DatabaseURL); err == nil {
			sqlDB = db
			app.DB = db
		} else {
			log.Printf("DB open failed: %v", err)
		}
	}
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		dbStatus := "disabled"
		if sqlDB != nil {
			if err := sqlDB.Ping(); err == nil {
				dbStatus = "ok"
			} else {
				dbStatus = "down"
			}
		}
		_, _ = w.Write([]byte(`{"status":"ok","deps":{"db":"` + dbStatus + `"}}`))
	})
	mux.HandleFunc("/api/auth/login", handlers.HandleLogin)
	// Protected routes
	mux.Handle("/api/conversations", middleware.JWT(http.HandlerFunc(handlers.HandleListConversations)))
	mux.Handle("/api/messages", middleware.JWT(http.HandlerFunc(handlers.HandleListMessages)))

	mux.HandleFunc("/api/roles", handlers.HandleRoles)
	mux.HandleFunc("/api/roles/", handlers.HandleRoleDetail)
	mux.HandleFunc("/api/chat", handlers.HandleChat)
	mux.HandleFunc("/api/chat/stream", handlers.HandleChatStream)
	mux.HandleFunc("/api/asr", handlers.HandleASR)
	mux.HandleFunc("/api/tts", handlers.HandleTTS)

	handler := withCORS(cfg.AllowedOrigins, middleware.RateLimit(100)(middleware.Logging(mux)))
	log.Printf("Go backend listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatal(err)
	}
}

func withCORS(allowedOrigins string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
