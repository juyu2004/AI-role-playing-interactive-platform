package memory

import (
	"ai-role-playing-platform/backend-go/internal/repo"
	"fmt"
	"sync"
	"time"
)

type ConversationRepo struct {
	mu    sync.RWMutex
	items []repo.Conversation
}

func NewConversationRepo() repo.ConversationRepository { return &ConversationRepo{} }

func (r *ConversationRepo) ListByUser(userID string, limit, offset int) ([]repo.Conversation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]repo.Conversation, 0)
	for _, c := range r.items {
		if c.UserID == userID {
			out = append(out, c)
		}
	}
	if offset > len(out) {
		return []repo.Conversation{}, nil
	}
	end := offset + limit
	if end > len(out) {
		end = len(out)
	}
	return out[offset:end], nil
}

func (r *ConversationRepo) Create(userID, roleID string) (repo.Conversation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := fmt.Sprintf("conv-%d", time.Now().UnixNano())
	c := repo.Conversation{ID: id, UserID: userID, RoleID: roleID}
	r.items = append(r.items, c)
	return c, nil
}

func (r *ConversationRepo) ReassignUserData(fromUserID, toUserID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.items {
		if r.items[i].UserID == fromUserID {
			r.items[i].UserID = toUserID
		}
	}
	return nil
}
