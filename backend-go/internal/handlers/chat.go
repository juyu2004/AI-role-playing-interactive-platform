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
	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		badRequest(w, "invalid json")
		return
	}

	// Load role strictly from DB for persistence
	var role *models.Role
	if r, err := rdb.NewRoleRepo(app.DB).GetByID(req.RoleID); err == nil {
		role = r
	}
	if role == nil {
		notFound(w, "role not found")
		return
	}
	// Build brief history context (last 10 messages)
	history := ""
	if req.ConversationID != "" && app.DB != nil {
		msgRepo := rdb.NewMessageRepo(app.DB)
		if msgs, err := msgRepo.ListByConversation(req.ConversationID, 10, 0); err == nil {
			// oldest to newest
			if len(msgs) > 0 {
				for i := len(msgs) - 1; i >= 0; i-- {
				}
				// but repo likely already returns newest first; rebuild in forward order
				for i := len(msgs) - 1; i >= 0; i-- {
					m := msgs[i]
					if m.Sender == "user" {
						history += "用户：" + m.Text + "\n"
					} else {
						history += "角色：" + m.Text + "\n"
					}
				}
			}
		}
	}
	// Strengthen system prompt: enforce language and persona
	systemPrompt := role.Prompt + "\n\n[指令]\n- 一律使用简体中文回答。\n- 避免无意义的前缀/问候，直接围绕用户问题作答。\n- 回答不超过 3 段，每段不超过 2 句。\n"
	userText := req.Text
	if history != "" {
		userText = "对话历史（供参考，不要重复）：\n" + history + "\n当前用户：" + req.Text
	}
	router := services.ProviderRouter{}
	provider := router.ResolveProvider(role.ID)
	reply, _ := provider.Generate(systemPrompt, userText)
	resp := models.ChatResponse{Text: reply, AudioURL: nil}
	// Persist if authenticated and DB configured; reuse or create conversation
	if app.DB != nil {
		if user, _ := r.Context().Value("user").(string); user != "" {
			convRepo := rdb.NewConversationRepo(app.DB)
			msgRepo := rdb.NewMessageRepo(app.DB)
			var convID string
			if req.ConversationID != "" {
				convID = req.ConversationID
			} else {
				conv, _ := convRepo.Create(user, req.RoleID)
				convID = conv.ID
			}
			_, _ = msgRepo.Create(convID, "user", req.Text)
			_, _ = msgRepo.Create(convID, "role", reply)
			resp.ConversationID = convID
		}
	}
	ok(w, resp)
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
		_, _ = w.Write([]byte("event: error\ndata: streaming not supported\n\n"))
		return
	}
	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_, _ = w.Write([]byte("event: error\ndata: invalid json\n\n"))
		flusher.Flush()
		return
	}
	var role *models.Role
	if rdbRole, err := rdb.NewRoleRepo(app.DB).GetByID(req.RoleID); err == nil {
		role = rdbRole
	}
	if role == nil {
		_, _ = w.Write([]byte("event: error\ndata: role not found\n\n"))
		flusher.Flush()
		return
	}
	// Build brief history
	history := ""
	if req.ConversationID != "" && app.DB != nil {
		msgRepo := rdb.NewMessageRepo(app.DB)
		if msgs, err := msgRepo.ListByConversation(req.ConversationID, 10, 0); err == nil {
			for i := len(msgs) - 1; i >= 0; i-- {
				m := msgs[i]
				if m.Sender == "user" {
					history += "用户：" + m.Text + "\n"
				} else {
					history += "角色：" + m.Text + "\n"
				}
			}
		}
	}
	systemPrompt := role.Prompt + "\n\n[指令]\n- 一律使用简体中文回答。\n- 避免无意义的前缀/问候，直接围绕用户问题作答。\n- 回答不超过 3 段，每段不超过 2 句。\n"
	userText := req.Text
	if history != "" {
		userText = "对话历史（供参考，不要重复）：\n" + history + "\n当前用户：" + req.Text
	}
	router := services.ProviderRouter{}
	provider := router.ResolveProvider(role.ID)
	reply, err := provider.Generate(systemPrompt, userText)
	if err != nil {
		_, _ = w.Write([]byte("event: error\ndata: generation failed\n\n"))
		flusher.Flush()
		return
	}
	const chunkSize = 120
	for i := 0; i < len(reply); i += chunkSize {
		end := i + chunkSize
		if end > len(reply) {
			end = len(reply)
		}
		chunk := reply[i:end]
		fmt.Fprintf(w, "event: message\n")
		fmt.Fprintf(w, "data: %s\n\n", escapeForSSE(chunk))
		flusher.Flush()
		time.Sleep(30 * time.Millisecond)
	}
	_, _ = w.Write([]byte("event: done\ndata: end\n\n"))
	flusher.Flush()
}

func escapeForSSE(s string) string {
	b, _ := json.Marshal(s)
	if len(b) >= 2 && b[0] == '"' && b[len(b)-1] == '"' {
		return string(b[1 : len(b)-1])
	}
	return string(b)
}
