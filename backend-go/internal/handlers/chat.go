package handlers

import (
    "encoding/json"
    "net/http"

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
    reply := services.GenerateReply(role.Prompt, req.Text)
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


