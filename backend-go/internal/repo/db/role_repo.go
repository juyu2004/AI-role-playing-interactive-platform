package db

import (
	"ai-role-playing-platform/backend-go/internal/models"
	"ai-role-playing-platform/backend-go/internal/repo"
	"database/sql"
	"log"
)

type RoleRepo struct{ db *sql.DB }

func NewRoleRepo(db *sql.DB) repo.RoleRepository { return &RoleRepo{db: db} }

func (r *RoleRepo) List() ([]models.Role, error) {
	log.Printf("db: roles List")
	rows, err := r.db.Query(`SELECT id, name, category, avatar_url, description, prompt FROM roles ORDER BY name`)
	if err != nil {
		log.Printf("db error: roles List: %v", err)
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
	log.Printf("db: roles GetByID id=%s", id)
	var m models.Role
	err := r.db.QueryRow(`SELECT id, name, category, avatar_url, description, prompt FROM roles WHERE id=$1`, id).
		Scan(&m.ID, &m.Name, &m.Category, &m.AvatarURL, &m.Description, &m.Prompt)
	if err != nil {
		log.Printf("db error: roles GetByID: %v", err)
		return nil, err
	}
	return &m, nil
}

func (r *RoleRepo) Create(role models.Role) (*models.Role, error) {
	log.Printf("db: roles Create id=%s name=%s", role.ID, role.Name)
	_, err := r.db.Exec(`INSERT INTO roles(id, name, category, avatar_url, description, prompt, image_url) VALUES($1,$2,$3,$4,$5,$6,$7)`,
		role.ID, role.Name, role.Category, role.AvatarURL, role.Description, role.Prompt, nil)
	if err != nil {
		log.Printf("db error: roles Create: %v", err)
		return nil, err
	}
	return r.GetByID(role.ID)
}

func (r *RoleRepo) UpdatePrompt(id string, prompt string) error {
	log.Printf("db: roles UpdatePrompt id=%s", id)
	_, err := r.db.Exec(`UPDATE roles SET prompt=$1 WHERE id=$2`, prompt, id)
	return err
}
