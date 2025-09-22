package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"ai-role-playing-platform/backend-go/internal/config"
	"ai-role-playing-platform/backend-go/internal/repo"
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
	_ = json.NewDecoder(r.Body).Decode(&req)
	// Login-or-register using in-memory user repo
	var userRepo repo.UserRepository = mem.NewUserRepo()
	u, err := userRepo.GetByEmail(req.Email)
	if err != nil {
		// register
		hash, _ := mem.HashPassword(req.Password)
		_ = userRepo.Create(&repo.User{ID: req.Email, Email: req.Email, PasswordHash: hash})
		u, _ = userRepo.GetByEmail(req.Email)
	} else {
		if !mem.CheckPassword(u.PasswordHash, req.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
			return
		}
	}
	cfg := config.Load()
	token, _ := security.SignHS256(cfg.JWTSecret, security.Claims{Sub: req.Email, Exp: time.Now().Add(24 * time.Hour).Unix()})
	_ = json.NewEncoder(w).Encode(loginResponse{Token: token})
}
