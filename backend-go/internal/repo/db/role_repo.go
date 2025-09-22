package db

import (
	"ai-role-playing-platform/backend-go/internal/models"
	"ai-role-playing-platform/backend-go/internal/repo"
	"database/sql"
)

type RoleRepo struct{ db *sql.DB }

func NewRoleRepo(db *sql.DB) repo.RoleRepository { return &RoleRepo{db: db} }

func (r *RoleRepo) List() ([]models.Role, error) {
	rows, err := r.db.Query(`SELECT id, name, category, avatar_url, description, prompt FROM roles ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Role
	for rows.Next() {
		var m models.Role
		if err := rows.Scan(&m.ID, &m.Name, &m.Category, &m.AvatarURL, &m.Description, &m.Prompt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *RoleRepo) GetByID(id string) (*models.Role, error) {
	var m models.Role
	err := r.db.QueryRow(`SELECT id, name, category, avatar_url, description, prompt FROM roles WHERE id=$1`, id).
		Scan(&m.ID, &m.Name, &m.Category, &m.AvatarURL, &m.Description, &m.Prompt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
