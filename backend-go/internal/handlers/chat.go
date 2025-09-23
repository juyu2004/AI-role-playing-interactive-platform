package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"ai-role-playing-platform/backend-go/internal/app"
	"ai-role-playing-platform/backend-go/internal/models"
	rdb "ai-role-playing-platform/backend-go/internal/repo/db"
	"ai-role-playing-platform/backend-go/internal/services"
)

func HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid json"})
		return
	}

	// Load role strictly from DB for persistence
	var role *models.Role
	if r, err := rdb.NewRoleRepo(app.DB).GetByID(req.RoleID); err == nil {
		role = r
	}
	if role == nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "role not found"})
		return
	}
	router := services.ProviderRouter{}
	provider := router.ResolveProvider(role.ID)
	reply, _ := provider.Generate(role.Prompt, req.Text)
	resp := models.ChatResponse{Text: reply, AudioURL: nil}
	// Persist if authenticated and DB configured (reuse global app.DB)
	if app.DB != nil {
		if user, _ := r.Context().Value("user").(string); user != "" {
			convRepo := rdb.NewConversationRepo(app.DB)
			msgRepo := rdb.NewMessageRepo(app.DB)
			conv, _ := convRepo.Create(user, req.RoleID)
			_, _ = msgRepo.Create(conv.ID, "user", req.Text)
			_, _ = msgRepo.Create(conv.ID, "role", reply)
		}
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// SSE placeholder: emits 3 chunks then closes.
func HandleChatStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for i := 1; i <= 3; i++ {
		fmt.Fprintf(w, "data: %s\n\n", fmt.Sprintf("chunk-%d", i))
		flusher.Flush()
		time.Sleep(300 * time.Millisecond)
	}
}
