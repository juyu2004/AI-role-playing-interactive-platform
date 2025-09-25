package repo

import "ai-role-playing-platform/backend-go/internal/models"

type RoleRepository interface {
	List() ([]models.Role, error)
	GetByID(id string) (*models.Role, error)
	Create(role models.Role) (*models.Role, error)
	UpdatePrompt(id string, prompt string) error
}

type User struct {
	ID           string
	Email        string
	PasswordHash string
	AvatarURL    *string
	Bio          *string
}

type UserRepository interface {
	Create(u *User) error
	GetByEmail(email string) (*User, error)
	GetByID(id string) (*User, error)
	UpdateAvatar(userID string, avatarURL *string) error
	UpdateBio(userID string, bio *string) error
}

type Conversation struct {
	ID     string
	UserID string
	RoleID string
}

type Message struct {
	ID             string
	ConversationID string
	Sender         string // user/role
	Text           string
}

type ConversationRepository interface {
	ListByUser(userID string, limit, offset int) ([]Conversation, error)
	Create(userID, roleID string) (Conversation, error)
	ReassignUserData(fromUserID, toUserID string) error
}

type MessageRepository interface {
	ListByConversation(convID string, limit, offset int) ([]Message, error)
	Create(convID, sender, text string) (Message, error)
}
