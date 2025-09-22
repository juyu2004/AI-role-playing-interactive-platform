package handlers

import (
	"encoding/json"
	"net/http"

	"ai-role-playing-platform/backend-go/internal/app"
	rdb "ai-role-playing-platform/backend-go/internal/repo/db"
)

type listMessagesResponse struct {
	Items []any `json:"items"`
}

// GET /api/messages (stub) ?conversationId=...
func HandleListMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if app.DB == nil {
		_ = json.NewEncoder(w).Encode(listMessagesResponse{Items: []any{}})
		return
	}
	convID := r.URL.Query().Get("conversationId")
	items, _ := rdb.NewMessageRepo(app.DB).ListByConversation(convID, 20, 0)
	_ = json.NewEncoder(w).Encode(listMessagesResponse{Items: toAnySlice(items)})
}
