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
	repo := rdb.NewRoleRepo(app.DB)
	roles, err := repo.List()
	if err != nil {
		serverError(w, "failed to list roles")
		return
	}
	ok(w, roles)
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
		notFound(w, "role not found")
		return
	}
	ok(w, role)
}

// POST /api/roles
func HandleCreateRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if app.DB == nil {
		serverError(w, "db not available")
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
		badRequest(w, "invalid json")
		return
	}
	// Use repository create
	repo := rdb.NewRoleRepo(app.DB)
	createdRole, err := repo.Create(models.Role{ID: req.ID, Name: req.Name, Category: req.Category, AvatarURL: derefOrEmpty(req.AvatarURL), Description: req.Description, Prompt: req.Prompt})
	if err != nil {
		serverError(w, "failed to create role")
		return
	}
	created(w, createdRole)
}

func derefOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
