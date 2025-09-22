package memory

import (
	"ai-role-playing-platform/backend-go/internal/repo"
	"fmt"
	"sync"
	"time"
)

type MessageRepo struct {
	mu    sync.RWMutex
	items []repo.Message
}

func NewMessageRepo() repo.MessageRepository { return &MessageRepo{} }

func (r *MessageRepo) ListByConversation(convID string, limit, offset int) ([]repo.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]repo.Message, 0)
	for _, m := range r.items {
		if m.ConversationID == convID {
			out = append(out, m)
		}
	}
	if offset > len(out) {
		return []repo.Message{}, nil
	}
	end := offset + limit
	if end > len(out) {
		end = len(out)
	}
	return out[offset:end], nil
}

func (r *MessageRepo) Create(convID, sender, text string) (repo.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := fmt.Sprintf("msg-%d", time.Now().UnixNano())
	m := repo.Message{ID: id, ConversationID: convID, Sender: sender, Text: text}
	r.items = append(r.items, m)
	return m, nil
}
