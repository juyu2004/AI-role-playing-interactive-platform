package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"ai-role-playing-platform/backend-go/internal/models"
	"ai-role-playing-platform/backend-go/internal/repo"
	mem "ai-role-playing-platform/backend-go/internal/repo/memory"
)

var inMemoryRoles = []models.Role{
	{
		ID:          "sherlock",
		Name:        "Sherlock Holmes",
		Category:    "literature",
		AvatarURL:   "",
		Description: "Brilliant detective with keen observation and deduction.",
		Prompt:      "You are Sherlock Holmes. Respond concisely with sharp deductions.",
	},
	{
		ID:          "mulan",
		Name:        "Hua Mulan",
		Category:    "history",
		AvatarURL:   "",
		Description: "Heroine known for courage and loyalty.",
		Prompt:      "You are Hua Mulan. Speak with bravery and humility.",
	},
}

var roleRepository repo.RoleRepository = mem.NewRoleRepo(inMemoryRoles)

func HandleRoles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	roles, _ := roleRepository.List()
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
	role, err := roleRepository.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(role)
}
