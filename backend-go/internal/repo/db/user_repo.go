package db

import (
	"ai-role-playing-platform/backend-go/internal/repo"
	"database/sql"
	"errors"
	"log"
)

type UserRepo struct{ db *sql.DB }

func NewUserRepo(db *sql.DB) repo.UserRepository { return &UserRepo{db: db} }

func (r *UserRepo) Create(u *repo.User) error {
	// Insert and return generated UUID id
	if u == nil {
		return errors.New("nil user")
	}
	log.Printf("db: users Create email=%s", u.Email)
	return r.db.QueryRow(`INSERT INTO users(email, password_hash) VALUES($1,$2) RETURNING id`, u.Email, u.PasswordHash).Scan(&u.ID)
}

func (r *UserRepo) GetByEmail(email string) (*repo.User, error) {
	log.Printf("db: users GetByEmail email=%s", email)
	var u repo.User
	err := r.db.QueryRow(`SELECT id, email, password_hash, avatar_url, bio FROM users WHERE email=$1`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.AvatarURL, &u.Bio)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByID(id string) (*repo.User, error) {
	log.Printf("db: users GetByID id=%s", id)
	var u repo.User
	err := r.db.QueryRow(`SELECT id, email, password_hash, avatar_url, bio FROM users WHERE id=$1`, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.AvatarURL, &u.Bio)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) UpdateAvatar(userID string, avatarURL *string) error {
	log.Printf("db: users UpdateAvatar id=%s", userID)
	_, err := r.db.Exec(`UPDATE users SET avatar_url=$1 WHERE id=$2`, avatarURL, userID)
	return err
}

func (r *UserRepo) UpdateBio(userID string, bio *string) error {
	log.Printf("db: users UpdateBio id=%s", userID)
	_, err := r.db.Exec(`UPDATE users SET bio=$1 WHERE id=$2`, bio, userID)
	return err
}
