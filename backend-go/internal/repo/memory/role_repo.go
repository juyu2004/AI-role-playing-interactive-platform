package memory

import (
	"ai-role-playing-platform/backend-go/internal/models"
	"ai-role-playing-platform/backend-go/internal/repo"
	"errors"
)

type RoleRepo struct {
	items map[string]models.Role
}

func NewRoleRepo(seed []models.Role) repo.RoleRepository {
	m := make(map[string]models.Role)
	for _, r := range seed {
		m[r.ID] = r
	}
	return &RoleRepo{items: m}
}

func (r *RoleRepo) List() ([]models.Role, error) {
	out := make([]models.Role, 0, len(r.items))
	for _, v := range r.items {
		out = append(out, v)
	}
	return out, nil
}

func (r *RoleRepo) GetByID(id string) (*models.Role, error) {
	v, ok := r.items[id]
	if !ok {
		return nil, errors.New("not found")
	}
	vv := v
	return &vv, nil
}

func (r *RoleRepo) Create(role models.Role) (*models.Role, error) {
    if _, exists := r.items[role.ID]; exists {
        return nil, errors.New("exists")
    }
    r.items[role.ID] = role
    rr := role
    return &rr, nil
}