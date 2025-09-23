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

// POST /api/conversations/reassign {"toUserId":"..."}
func HandleReassignUserData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if app.DB == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "db not configured"})
		return
	}
	var req struct {
		ToUserID string `json:"toUserId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ToUserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fromUserID, _ := r.Context().Value("user").(string)
	if fromUserID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := rdb.NewConversationRepo(app.DB).ReassignUserData(fromUserID, req.ToUserID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "reassign failed"})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func toAnySlice[T any](xs []T) []any {
	out := make([]any, len(xs))
	for i, v := range xs {
		out[i] = v
	}
	return out
}
