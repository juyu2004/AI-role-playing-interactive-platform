package memory

import "ai-role-playing-platform/backend-go/internal/repo"

type MessageRepo struct{}

func NewMessageRepo() repo.MessageRepository { return &MessageRepo{} }

func (r *MessageRepo) ListByConversation(convID string, limit, offset int) ([]repo.Message, error) {
	return []repo.Message{}, nil
}
