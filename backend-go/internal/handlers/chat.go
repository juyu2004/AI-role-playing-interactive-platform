package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"ai-role-playing-platform/backend-go/internal/models"
	"ai-role-playing-platform/backend-go/internal/services"
)

func HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role := getRoleByID(req.RoleID)
	if role == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	router := services.ProviderRouter{}
	provider := router.ResolveProvider(role.ID)
	reply, _ := provider.Generate(role.Prompt, req.Text)
	resp := models.ChatResponse{Text: reply, AudioURL: nil}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func getRoleByID(id string) *models.Role {
	for _, role := range inMemoryRoles {
		if role.ID == id {
			r := role
			return &r
		}
	}
	return nil
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
