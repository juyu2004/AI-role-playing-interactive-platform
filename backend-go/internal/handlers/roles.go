package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"ai-role-playing-platform/backend-go/internal/app"
	"ai-role-playing-platform/backend-go/internal/models"
	rdb "ai-role-playing-platform/backend-go/internal/repo/db"
)

// Roles are stored in PostgreSQL; seeded by migrations.

func HandleRoles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	repo := rdb.NewRoleRepo(app.DB)
	roles, _ := repo.List()
	_ = json.NewEncoder(w).Encode(roles)
}

func HandleRoleDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// path: /api/roles/{id}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/roles/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := parts[0]
	repo := rdb.NewRoleRepo(app.DB)
	role, err := repo.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(role)
}

// POST /api/roles
func HandleCreateRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if app.DB == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	var req struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Category    string  `json:"category"`
		AvatarURL   *string `json:"avatarUrl"`
		Description string  `json:"description"`
		Prompt      string  `json:"prompt"`
		ImageURL    *string `json:"imageUrl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Use repository create
	repo := rdb.NewRoleRepo(app.DB)
	created, err := repo.Create(models.Role{ID: req.ID, Name: req.Name, Category: req.Category, AvatarURL: derefOrEmpty(req.AvatarURL), Description: req.Description, Prompt: req.Prompt})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(created)
}

func derefOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
