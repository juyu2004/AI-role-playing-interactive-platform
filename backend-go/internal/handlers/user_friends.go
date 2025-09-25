package handlers

import (
	"encoding/json"
	"net/http"

	"ai-role-playing-platform/backend-go/internal/app"
	"strings"
)

// Simple friends endpoints using direct SQL against friends table

// GET /api/user/friends
func HandleListFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	userID, _ := r.Context().Value("user").(string)
	rows, err := app.DB.Query(`SELECT friend_user_id FROM friends WHERE user_id=$1`, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		_ = rows.Scan(&id)
		ids = append(ids, id)
	}
	ok(w, map[string]any{"friends": ids})
}

// POST /api/user/friends {"friendUserId":"..."}
func HandleAddFriend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	userID, _ := r.Context().Value("user").(string)
	var req struct {
		FriendUserID string `json:"friendUserId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.FriendUserID == "" {
		badRequest(w, "invalid json")
		return
	}
	// insert both directions for undirected
	if _, err := app.DB.Exec(`INSERT INTO friends(user_id, friend_user_id) VALUES($1,$2) ON CONFLICT DO NOTHING`, userID, req.FriendUserID); err != nil {
		serverError(w, "db error")
		return
	}
	if _, err := app.DB.Exec(`INSERT INTO friends(user_id, friend_user_id) VALUES($1,$2) ON CONFLICT DO NOTHING`, req.FriendUserID, userID); err != nil {
		serverError(w, "db error")
		return
	}
	ok(w, map[string]string{"status": "ok"})
}

// DELETE /api/user/friends/{friendUserId}
func HandleRemoveFriend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	userID, _ := r.Context().Value("user").(string)
	friendID := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	if friendID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = app.DB.Exec(`DELETE FROM friends WHERE (user_id=$1 AND friend_user_id=$2) OR (user_id=$2 AND friend_user_id=$1)`, userID, friendID)
	ok(w, map[string]string{"status": "ok"})
}
