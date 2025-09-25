package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"ai-role-playing-platform/backend-go/internal/app"
	"ai-role-playing-platform/backend-go/internal/config"
	"ai-role-playing-platform/backend-go/internal/repo"
	rdb "ai-role-playing-platform/backend-go/internal/repo/db"
	mem "ai-role-playing-platform/backend-go/internal/repo/memory"
	"ai-role-playing-platform/backend-go/internal/security"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// POST /api/auth/login (stub)
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	// Login-or-register using PostgreSQL user repo
	userRepo := rdb.NewUserRepo(app.DB)
	u, err := userRepo.GetByEmail(req.Email)
	if err != nil {
		hash, _ := mem.HashPassword(req.Password)
		_ = userRepo.Create(&repo.User{Email: req.Email, PasswordHash: hash})
		u, _ = userRepo.GetByEmail(req.Email)
	} else {
		if !mem.CheckPassword(u.PasswordHash, req.Password) {
			unauthorized(w, "invalid credentials")
			return
		}
	}
	cfg := config.Load()
	// Use user ID as subject
	token, _ := security.SignHS256(cfg.JWTSecret, security.Claims{Sub: u.ID, Exp: time.Now().Add(24 * time.Hour).Unix()})
	ok(w, loginResponse{Token: token})
}

// POST /api/auth/register
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	if req.Email == "" || req.Password == "" {
		badRequest(w, "email and password are required")
		return
	}
	userRepo := rdb.NewUserRepo(app.DB)
	if _, err := userRepo.GetByEmail(req.Email); err == nil {
		badRequest(w, "email already registered")
		return
	}
	hash, err := mem.HashPassword(req.Password)
	if err != nil {
		serverError(w, "failed to hash password")
		return
	}
	if err := userRepo.Create(&repo.User{Email: req.Email, PasswordHash: hash}); err != nil {
		serverError(w, "failed to create user")
		return
	}
	u, err := userRepo.GetByEmail(req.Email)
	if err != nil {
		serverError(w, "failed to load user")
		return
	}
	cfg := config.Load()
	token, _ := security.SignHS256(cfg.JWTSecret, security.Claims{Sub: u.ID, Exp: time.Now().Add(24 * time.Hour).Unix()})
	created(w, loginResponse{Token: token})
}
