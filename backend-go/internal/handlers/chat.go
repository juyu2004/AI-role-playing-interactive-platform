package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
	"unicode"

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
	systemPrompt := role.Prompt + "\n\n[指令]\n- 一律使用简体中文回答。\n- 严禁自我介绍/重复问候（如“我是" + role.Name + "，很高兴…”），除非用户明确询问“你是谁”。\n- 直接围绕用户问题作答，避免铺垫与客套。\n- 回答不超过 3 段，每段不超过 2 句。\n"
	userText := req.Text
	if history != "" {
		userText = "对话历史（供参考，不要重复）：\n" + history + "\n当前用户：" + req.Text
	}
	router := services.ProviderRouter{}
	provider := router.ResolveProvider(role.ID)
	reply, err := provider.Generate(systemPrompt, userText)
	if err != nil {
		// 如果API调用失败，返回错误信息
		reply = "抱歉,我暂时无法回答这个问题。"
		fmt.Printf("API调用失败: %v\n", err)
	} else {
		reply = sanitizeChinese(stripGreetings(reply, role.Name))
	}
	resp := models.ChatResponse{Text: reply, AudioURL: nil}
	// Persist if authenticated and DB configured; reuse latest conversation by user+role when未提供 conversationId
	if app.DB != nil {
		if user, _ := r.Context().Value("user").(string); user != "" {
			convRepo := rdb.NewConversationRepo(app.DB)
			msgRepo := rdb.NewMessageRepo(app.DB)
			var convID string
			if req.ConversationID != "" {
				convID = req.ConversationID
			} else {
				if existing, err := convRepo.GetLatestByUserRole(user, req.RoleID); err == nil && existing != nil {
					convID = existing.ID
				} else {
					conv, _ := convRepo.Create(user, req.RoleID)
					convID = conv.ID
				}
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
	systemPrompt := role.Prompt + "\n\n[指令]\n- 一律使用简体中文回答。\n- 严禁自我介绍/重复问候（如“我是" + role.Name + "，很高兴…”），除非用户明确询问“你是谁”。\n- 直接围绕用户问题作答，避免铺垫与客套。\n- 回答不超过 3 段，每段不超过 2 句。\n"
	userText := req.Text
	if history != "" {
		userText = "对话历史（供参考，不要重复）：\n" + history + "\n当前用户：" + req.Text
	}
	router := services.ProviderRouter{}
	provider := router.ResolveProvider(role.ID)
	reply, err := provider.Generate(systemPrompt, userText)
	if err != nil {
		// 如果API调用失败，返回错误信息
		reply = "抱歉,我暂时无法回答这个问题。"
		fmt.Printf("API调用失败: %v\n", err)
	} else {
		reply = sanitizeChinese(stripGreetings(reply, role.Name))
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

// sanitizeChinese 过滤英文代码感强的片段与异常符号，尽量输出中文与常用标点
var (
	reAlphaSeq = regexp.MustCompile(`[A-Za-z_]{3,}`)
	reWeird    = regexp.MustCompile("[`~@#$%^*_=<>\\{}\\[\\]|]+")
)

func sanitizeChinese(s string) string {
	// 如果中文占比很低，先删除长英文串与怪异符号
	total := 0
	han := 0
	for _, r := range s {
		total++
		if unicode.Is(unicode.Han, r) {
			han++
		}
	}
	cleaned := s
	if total == 0 {
		return s
	}
	ratio := float64(han) / float64(total)
	if ratio < 0.7 {
		cleaned = reAlphaSeq.ReplaceAllString(cleaned, "")
		cleaned = reWeird.ReplaceAllString(cleaned, "")
	}
	// 常见伪代码/占位词清理
	tokens := []string{"return", "args", "NIL", "epoch", "Talk", "cups"}
	for _, t := range tokens {
		cleaned = regexp.MustCompile(`(?i)`+regexp.QuoteMeta(t)).ReplaceAllString(cleaned, "")
	}
	// 收尾：去掉多余空格
	// 简化连续空格
	cleaned = regexp.MustCompile(`\s{2,}`).ReplaceAllString(cleaned, " ")
	return cleaned
}

// stripGreetings 去除模型常见的自我介绍与客套话
func stripGreetings(s, roleName string) string {
	patterns := []string{
		`^\s*我是[^，。,!！]*[，。,!！]\s*`,
		`^\s*我叫[^，。,!！]*[，。,!！]\s*`,
		`^\s*很高兴[^。.!！]*[。.!！]\s*`,
	}
	if roleName != "" {
		patterns = append(patterns, `^\s*我是\s*`+regexp.QuoteMeta(roleName)+`[^。.!！]*[。.!！]\s*`)
	}
	out := s
	for _, p := range patterns {
		out = regexp.MustCompile(p).ReplaceAllString(out, "")
	}
	// 句首若仍为空客套字词，裁掉
	out = regexp.MustCompile(`^(大家好|您好|你好)[，,\s]*`).ReplaceAllString(out, "")
	return out
}
