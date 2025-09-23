package db

import (
	"ai-role-playing-platform/backend-go/internal/repo"
	"database/sql"
	"log"
)

type ConversationRepo struct{ db *sql.DB }

func NewConversationRepo(db *sql.DB) repo.ConversationRepository { return &ConversationRepo{db: db} }

func (r *ConversationRepo) ListByUser(userID string, limit, offset int) ([]repo.Conversation, error) {
	log.Printf("db: conversations ListByUser userID=%s limit=%d offset=%d", userID, limit, offset)
	rows, err := r.db.Query(`SELECT id, user_id, role_id FROM conversations WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		log.Printf("db error: conversations ListByUser: %v", err)
		return nil, err
	}
	defer rows.Close()
	var out []repo.Conversation
	for rows.Next() {
		var c repo.Conversation
		if err := rows.Scan(&c.ID, &c.UserID, &c.RoleID); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *ConversationRepo) Create(userID, roleID string) (repo.Conversation, error) {
	log.Printf("db: conversations Create userID=%s roleID=%s", userID, roleID)
	var id string
	err := r.db.QueryRow(`INSERT INTO conversations(user_id, role_id) VALUES($1,$2) RETURNING id`, userID, roleID).Scan(&id)
	if err != nil {
		log.Printf("db error: conversations Create: %v", err)
		return repo.Conversation{}, err
	}
	return repo.Conversation{ID: id, UserID: userID, RoleID: roleID}, nil
}

func (r *ConversationRepo) ReassignUserData(fromUserID, toUserID string) error {
	log.Printf("db: conversations ReassignUserData from=%s to=%s", fromUserID, toUserID)
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.Exec(`UPDATE conversations SET user_id=$1 WHERE user_id=$2`, toUserID, fromUserID); err != nil {
		log.Printf("db error: conversations ReassignUserData: %v", err)
		return err
	}
	return tx.Commit()
}
