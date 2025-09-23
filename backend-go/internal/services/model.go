package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"ai-role-playing-platform/backend-go/internal/config"
)

type ModelProvider interface {
	Generate(personaPrompt string, userText string) (string, error)
}

type MockProvider struct{}

func (m MockProvider) Generate(personaPrompt string, userText string) (string, error) {
	if personaPrompt == "" {
		return fmt.Sprintf("(mock) You said: %s", userText), nil
	}
	return fmt.Sprintf("(mock persona) %s | Replying to: %s", personaPrompt, userText), nil
}

// Router chooses provider by role/model mapping. Placeholder returns MockProvider.
type ProviderRouter struct{}

func (r ProviderRouter) ResolveProvider(roleID string) ModelProvider {
	// If OpenRouter is configured, use it; otherwise fallback to mock
	cfg := config.Load()
	if cfg.OpenRouterAPIKey != "" {
		return OpenRouterProvider{
			APIKey:  cfg.OpenRouterAPIKey,
			BaseURL: cfg.OpenRouterBaseURL,
			Model:   cfg.OpenRouterModel,
			Client:  &http.Client{Timeout: 30 * time.Second},
		}
	}
	return MockProvider{}
}

// OpenRouterProvider implements ModelProvider using OpenRouter chat/completions API
type OpenRouterProvider struct {
	APIKey  string
	BaseURL string
	Model   string
	Client  *http.Client
}

type orMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type orRequest struct {
	Model    string      `json:"model"`
	Messages []orMessage `json:"messages"`
}

type orChoiceMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type orChoice struct {
	Message orChoiceMessage `json:"message"`
}

type orResponse struct {
	Choices []orChoice `json:"choices"`
}

func (p OpenRouterProvider) Generate(personaPrompt string, userText string) (string, error) {
	if p.Client == nil {
		p.Client = &http.Client{Timeout: 30 * time.Second}
	}
	reqBody := orRequest{
		Model: p.Model,
		Messages: []orMessage{
			{Role: "system", Content: personaPrompt},
			{Role: "user", Content: userText},
		},
	}
	b, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("%s/chat/completions", p.BaseURL)
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.APIKey)
	// Optional: identify app per OpenRouter guidelines
	httpReq.Header.Set("HTTP-Referer", "https://localhost")
	httpReq.Header.Set("X-Title", "AI Role Playing Platform")

	resp, err := p.Client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("openrouter http %d", resp.StatusCode)
	}
	var orResp orResponse
	if err := json.NewDecoder(resp.Body).Decode(&orResp); err != nil {
		return "", err
	}
	if len(orResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}
	return orResp.Choices[0].Message.Content, nil
}
