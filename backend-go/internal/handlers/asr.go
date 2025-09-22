package handlers

import (
	"encoding/json"
	"net/http"

	"ai-role-playing-platform/backend-go/internal/services"
)

type asrResponse struct {
	Text string `json:"text"`
}

// POST /api/asr  (stub)
func HandleASR(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// In the future: parse audio from multipart/form-data or raw bytes
	text, _ := services.TranscribeAudio(nil)
	_ = json.NewEncoder(w).Encode(asrResponse{Text: text})
}
