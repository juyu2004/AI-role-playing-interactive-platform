package db

import (
	"ai-role-playing-platform/backend-go/internal/repo"
	"database/sql"
	"log"
)

type MessageRepo struct{ db *sql.DB }

func NewMessageRepo(db *sql.DB) repo.MessageRepository { return &MessageRepo{db: db} }

func (r *MessageRepo) ListByConversation(convID string, limit, offset int) ([]repo.Message, error) {
	log.Printf("db: messages ListByConversation convID=%s limit=%d offset=%d", convID, limit, offset)
	rows, err := r.db.Query(`SELECT id, conversation_id, sender, text FROM messages WHERE conversation_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, convID, limit, offset)
	if err != nil {
		log.Printf("db error: messages ListByConversation: %v", err)
		return nil, err
	}
	defer rows.Close()
	var out []repo.Message
	for rows.Next() {
		var m repo.Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.Sender, &m.Text); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *MessageRepo) Create(convID, sender, text string) (repo.Message, error) {
	log.Printf("db: messages Create convID=%s sender=%s", convID, sender)
	var id string
	err := r.db.QueryRow(`INSERT INTO messages(conversation_id, sender, text) VALUES($1,$2,$3) RETURNING id`, convID, sender, text).Scan(&id)
	if err != nil {
		log.Printf("db error: messages Create: %v", err)
		return repo.Message{}, err
	}
	return repo.Message{ID: id, ConversationID: convID, Sender: sender, Text: text}, nil
}
