package models

type Role struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	AvatarURL    string `json:"avatarUrl"`
	Description  string `json:"description"`
	Prompt       string `json:"-"`
	SystemPrompt string `json:"-"`
	UserGuide    string `json:"-"`
}

type ChatRequest struct {
	RoleID         string `json:"roleId"`
	Text           string `json:"text"`
	ConversationID string `json:"conversationId,omitempty"`
}

type ChatResponse struct {
	Text           string  `json:"text"`
	AudioURL       *string `json:"audioUrl"`
	ConversationID string  `json:"conversationId,omitempty"`
}
