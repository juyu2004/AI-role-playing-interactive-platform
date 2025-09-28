package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"ai-role-playing-platform/backend-go/internal/app"
	"ai-role-playing-platform/backend-go/internal/config"
	"ai-role-playing-platform/backend-go/internal/handlers"
	"ai-role-playing-platform/backend-go/internal/middleware"
	"ai-role-playing-platform/backend-go/internal/models"
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
	mux.HandleFunc("/api/auth/register", handlers.HandleRegister)
	// Protected routes
	mux.Handle("/api/conversations", middleware.JWT(http.HandlerFunc(handlers.HandleListConversations)))
	mux.Handle("/api/conversations/reassign", middleware.JWT(http.HandlerFunc(handlers.HandleReassignUserData)))
	mux.Handle("/api/messages", middleware.JWT(http.HandlerFunc(handlers.HandleListMessages)))

	mux.HandleFunc("/api/roles", handlers.HandleRoles)
	mux.HandleFunc("/api/roles/", handlers.HandleRoleDetail)
	mux.HandleFunc("/api/roles/prompt/", handlers.HandleUpdateRolePrompt)
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
	mux.HandleFunc("/api/upload/audio", handlers.HandleUploadAudio)

	// Seed roles from JSON if available
	if sqlDB != nil {
		if err := seedRolesFromJSON(sqlDB); err != nil {
			log.Printf("role seed skipped: %v", err)
		}
	}

	handler := withCORS(cfg.AllowedOrigins, middleware.RateLimit(100)(middleware.Logging(mux)))
	log.Printf("Go backend listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatal(err)
	}
}

type seedRole struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Prompt   string `json:"prompt"`
}

type seedFile struct {
	Roles []seedRole `json:"roles"`
}

func seedRolesFromJSON(db *sql.DB) error {
	// Try ../role_prompt/prompt.json relative to backend-go working dir
	candidates := []string{"../role_prompt/prompt.json", "../../role_prompt/prompt.json", "role_prompt/prompt.json"}
	var f *os.File
	var err error
	for _, p := range candidates {
		if _, statErr := os.Stat(p); statErr == nil {
			f, err = os.Open(p)
			if err == nil {
				defer f.Close()
				var sf seedFile
				if decodeErr := json.NewDecoder(f).Decode(&sf); decodeErr != nil {
					return decodeErr
				}
				repo := rdb.NewRoleRepo(db)
				for _, r := range sf.Roles {
					id := slug(r.Name)
					if existing, getErr := repo.GetByID(id); getErr == nil && existing != nil {
						// Update only prompt
						_ = repo.UpdatePrompt(id, r.Prompt)
					} else {
						// Create new role
						_, _ = repo.Create(models.Role{ID: id, Name: r.Name, Category: r.Category, AvatarURL: "", Description: r.Name, Prompt: r.Prompt})
					}
				}
				log.Printf("seeded roles from %s", p)
				return nil
			}
		}
	}
	return nil
}

func slug(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	// Replace spaces and non-alnum with hyphen, keep basic unicode letters/digits
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			prevDash = false
		} else if r == ' ' || r == '_' || r == '-' {
			if !prevDash {
				b.WriteByte('-')
				prevDash = true
			}
		} else {
			// drop other chars
		}
	}
	out := b.String()
	if out == "" {
		out = "role"
	}
	return out
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
