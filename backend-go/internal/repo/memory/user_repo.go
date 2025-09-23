package memory

import (
	"ai-role-playing-platform/backend-go/internal/repo"
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	mu      sync.RWMutex
	byEmail map[string]*repo.User
	byID    map[string]*repo.User
}

func NewUserRepo() repo.UserRepository {
	return &UserRepo{byEmail: make(map[string]*repo.User), byID: make(map[string]*repo.User)}
}

func (r *UserRepo) Create(u *repo.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byEmail[u.Email]; ok {
		return errors.New("exists")
	}
	if u.ID == "" {
		u.ID = u.Email
	}
	r.byEmail[u.Email] = u
	r.byID[u.ID] = u
	return nil
}

func (r *UserRepo) GetByEmail(email string) (*repo.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if u, ok := r.byEmail[email]; ok {
		return u, nil
	}
	return nil, errors.New("not found")
}

func (r *UserRepo) GetByID(id string) (*repo.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if u, ok := r.byID[id]; ok {
		return u, nil
	}
	return nil, errors.New("not found")
}

func (r *UserRepo) UpdateAvatar(userID string, avatarURL *string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.byID[userID]; ok {
		u.AvatarURL = avatarURL
		return nil
	}
	return errors.New("not found")
}

func (r *UserRepo) UpdateBio(userID string, bio *string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.byID[userID]; ok {
		u.Bio = bio
		return nil
	}
	return errors.New("not found")
}

func HashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func CheckPassword(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}
