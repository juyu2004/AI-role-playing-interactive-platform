package handlers

import (
	"encoding/json"
	"net/http"

	"ai-role-playing-platform/backend-go/internal/app"
	rdb "ai-role-playing-platform/backend-go/internal/repo/db"
)

// GET /api/user/me
func HandleGetMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	userID, _ := r.Context().Value("user").(string)
	u, err := rdb.NewUserRepo(app.DB).GetByID(userID)
	if err != nil {
		notFound(w, "user not found")
		return
	}
	ok(w, map[string]any{
		"id":        u.ID,
		"email":     u.Email,
		"avatarUrl": u.AvatarURL,
		"bio":       u.Bio,
	})
}

// POST /api/user/avatar {"avatarUrl":"..."}
func HandleUpdateAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	userID, _ := r.Context().Value("user").(string)
	var req struct {
		AvatarURL *string `json:"avatarUrl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	if err := rdb.NewUserRepo(app.DB).UpdateAvatar(userID, req.AvatarURL); err != nil {
		serverError(w, "failed to update avatar")
		return
	}
	ok(w, map[string]string{"status": "ok"})
}

// POST /api/user/bio {"bio":"..."}
func HandleUpdateBio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	userID, _ := r.Context().Value("user").(string)
	var req struct {
		Bio *string `json:"bio"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		badRequest(w, "invalid json")
		return
	}
	if err := rdb.NewUserRepo(app.DB).UpdateBio(userID, req.Bio); err != nil {
		serverError(w, "failed to update bio")
		return
	}
	ok(w, map[string]string{"status": "ok"})
}
