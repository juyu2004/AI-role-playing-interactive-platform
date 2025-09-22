package handlers

import (
	"encoding/json"
	"net/http"
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
	_ = json.NewEncoder(w).Encode(listMessagesResponse{Items: []any{}})
}
