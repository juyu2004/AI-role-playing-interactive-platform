package handlers

import (
	"encoding/json"
	"io"
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
	// Accept raw audio bytes or multipart
	var audio []byte
	if r.Header.Get("Content-Type") == "application/octet-stream" {
		audio, _ = io.ReadAll(r.Body)
	} else {
		// best-effort for other content-types
		audio, _ = io.ReadAll(r.Body)
	}
	text, err := services.TranscribeAudio(audio)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	_ = json.NewEncoder(w).Encode(asrResponse{Text: text})
}
