package memory

import "ai-role-playing-platform/backend-go/internal/repo"

type ConversationRepo struct{}

func NewConversationRepo() repo.ConversationRepository { return &ConversationRepo{} }

func (r *ConversationRepo) ListByUser(userID string, limit, offset int) ([]repo.Conversation, error) {
	return []repo.Conversation{}, nil
}
