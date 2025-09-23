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
	"ai-role-playing-platform/backend-go/internal/repo/db/migrate"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()
	var sqlDB *sql.DB
	if db, err := rdb.Open(cfg.DatabaseURL); err == nil {
		sqlDB = db
		app.DB = db
		if err := migrate.Apply(sqlDB); err != nil {
			log.Fatalf("DB migration failed: %v", err)
		}
	} else {
		log.Fatalf("DB open failed: %v", err)
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
	mux.Handle("/api/conversations/reassign", middleware.JWT(http.HandlerFunc(handlers.HandleReassignUserData)))
	mux.Handle("/api/messages", middleware.JWT(http.HandlerFunc(handlers.HandleListMessages)))

	mux.HandleFunc("/api/roles", handlers.HandleRoles)
	mux.HandleFunc("/api/roles/", handlers.HandleRoleDetail)
	mux.Handle("/api/user/me", middleware.JWT(http.HandlerFunc(handlers.HandleGetMe)))
	mux.Handle("/api/user/avatar", middleware.JWT(http.HandlerFunc(handlers.HandleUpdateAvatar)))
	mux.Handle("/api/user/bio", middleware.JWT(http.HandlerFunc(handlers.HandleUpdateBio)))
	mux.Handle("/api/user/friends", middleware.JWT(http.HandlerFunc(handlers.HandleListFriends)))
	mux.Handle("/api/user/friends/add", middleware.JWT(http.HandlerFunc(handlers.HandleAddFriend)))
	mux.Handle("/api/user/friends/", middleware.JWT(http.HandlerFunc(handlers.HandleRemoveFriend)))
	mux.HandleFunc("/api/roles/create", handlers.HandleCreateRole)
	mux.Handle("/api/chat", middleware.JWT(http.HandlerFunc(handlers.HandleChat)))
	mux.Handle("/api/chat/stream", middleware.JWT(http.HandlerFunc(handlers.HandleChatStream)))
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
