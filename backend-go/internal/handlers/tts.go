package handlers

import (
	"encoding/json"
	"net/http"

	"ai-role-playing-platform/backend-go/internal/services"
)

type ttsRequest struct {
	Text string `json:"text"`
}

type ttsResponse struct {
	AudioURL string `json:"audioUrl"`
}

// POST /api/tts  (stub)
func HandleTTS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req ttsRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	url, _ := services.SynthesizeSpeech(req.Text)
	_ = json.NewEncoder(w).Encode(ttsResponse{AudioURL: url})
}
