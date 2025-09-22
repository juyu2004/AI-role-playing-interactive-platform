package handlers

import (
	"encoding/json"
	"net/http"
)

type listConversationsResponse struct {
	Items []any `json:"items"`
}

// GET /api/conversations (stub)
func HandleListConversations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_ = json.NewEncoder(w).Encode(listConversationsResponse{Items: []any{}})
}
