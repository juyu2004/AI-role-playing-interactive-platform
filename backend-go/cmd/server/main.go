package main

import (
	"log"
	"net/http"

	"ai-role-playing-platform/backend-go/internal/config"
	"ai-role-playing-platform/backend-go/internal/handlers"
	"ai-role-playing-platform/backend-go/internal/middleware"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/api/auth/login", handlers.HandleLogin)
	mux.HandleFunc("/api/conversations", handlers.HandleListConversations)
	mux.HandleFunc("/api/messages", handlers.HandleListMessages)

	mux.HandleFunc("/api/roles", handlers.HandleRoles)
	mux.HandleFunc("/api/roles/", handlers.HandleRoleDetail)
	mux.HandleFunc("/api/chat", handlers.HandleChat)
	mux.HandleFunc("/api/chat/stream", handlers.HandleChatStream)
	mux.HandleFunc("/api/asr", handlers.HandleASR)
	mux.HandleFunc("/api/tts", handlers.HandleTTS)

	handler := withCORS(middleware.Logging(mux))
	log.Printf("Go backend listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatal(err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For now allow all; can be narrowed to cfg.AllowedOrigins if passed in
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
