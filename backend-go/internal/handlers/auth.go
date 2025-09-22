package handlers

import (
	"encoding/json"
	"net/http"
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
	_ = json.NewEncoder(w).Encode(loginResponse{Token: "dummy-token"})
}
