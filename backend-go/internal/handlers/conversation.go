package handlers

import (
	"encoding/json"
	"net/http"

	"ai-role-playing-platform/backend-go/internal/app"
	rdb "ai-role-playing-platform/backend-go/internal/repo/db"
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
	if app.DB == nil {
		_ = json.NewEncoder(w).Encode(listConversationsResponse{Items: []any{}})
		return
	}
	user, _ := r.Context().Value("user").(string)
	items, _ := rdb.NewConversationRepo(app.DB).ListByUser(user, 20, 0)
	_ = json.NewEncoder(w).Encode(listConversationsResponse{Items: toAnySlice(items)})
}

func toAnySlice[T any](xs []T) []any {
	out := make([]any, len(xs))
	for i, v := range xs {
		out[i] = v
	}
	return out
}
