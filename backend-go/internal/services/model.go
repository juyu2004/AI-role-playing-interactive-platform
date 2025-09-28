package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
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
	cfg := config.Load()
	// 优先使用阿里云百炼
	if cfg.AliyunAPIKey != "" {
		return AliyunProvider{
			APIKey:  cfg.AliyunAPIKey,
			BaseURL: cfg.AliyunBaseURL,
			Model:   cfg.AliyunModel,
			Client:  &http.Client{Timeout: 30 * time.Second},
		}
	}
	// 如果阿里云未配置，使用OpenRouter
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
	Model       string      `json:"model"`
	Messages    []orMessage `json:"messages"`
	MaxTokens   int         `json:"max_tokens,omitempty"`
	Temperature float64     `json:"temperature,omitempty"`
	TopP        float64     `json:"top_p,omitempty"`
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
		// Robust transport with timeouts
		transport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 20 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
		}
		p.Client = &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		}
	}
	reqBody := orRequest{
		Model: p.Model,
		Messages: []orMessage{
			{Role: "system", Content: personaPrompt},
			{Role: "user", Content: userText},
		},
		MaxTokens:   512,
		Temperature: 0.7,
		TopP:        0.9,
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

	// Context timeout per request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	httpReq = httpReq.WithContext(ctx)

	// Retry on transient errors with exponential backoff
	resp, err := p.Client.Do(httpReq)
	if err != nil || resp.StatusCode >= 500 || resp.StatusCode == 429 {
		// Longer backoff for rate limiting
		backoffTime := 2 * time.Second
		if resp != nil && resp.StatusCode == 429 {
			backoffTime = 5 * time.Second
		}
		time.Sleep(backoffTime)
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		resp, err = p.Client.Do(httpReq)
		if err != nil {
			return "", err
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorMsg string
		switch resp.StatusCode {
		case 429:
			errorMsg = "rate limit exceeded (429)"
		case 401:
			errorMsg = "unauthorized (401) - check API key"
		case 402:
			errorMsg = "payment required (402) - insufficient credits"
		case 404:
			errorMsg = "model not found (404)"
		default:
			errorMsg = fmt.Sprintf("http error %d", resp.StatusCode)
		}
		return "", fmt.Errorf("openrouter %s", errorMsg)
	}
	var orResp orResponse
	if err := json.NewDecoder(resp.Body).Decode(&orResp); err != nil {
		return "", err
	}
	if len(orResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}
	out := orResp.Choices[0].Message.Content
	// Short-answer mitigation: retry once with larger budget
	if len(out) < 20 {
		reqBody.MaxTokens = 768
		reqBody.Temperature = 0.6
		b2, _ := json.Marshal(reqBody)
		httpReq2, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(b2))
		httpReq2.Header = httpReq.Header.Clone()
		ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel2()
		httpReq2 = httpReq2.WithContext(ctx2)
		if resp2, err2 := p.Client.Do(httpReq2); err2 == nil {
			defer resp2.Body.Close()
			var orResp2 orResponse
			if err := json.NewDecoder(resp2.Body).Decode(&orResp2); err == nil && len(orResp2.Choices) > 0 {
				if c := orResp2.Choices[0].Message.Content; len(c) >= len(out) {
					out = c
				}
			}
		}
	}
	return out, nil
}

// AliyunProvider implements ModelProvider using Alibaba Cloud Bailian API
type AliyunProvider struct {
	APIKey  string
	BaseURL string
	Model   string
	Client  *http.Client
}

type aliyunMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type aliyunRequest struct {
	Model       string          `json:"model"`
	Messages    []aliyunMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	TopP        float64         `json:"top_p,omitempty"`
}

type aliyunChoiceMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type aliyunChoice struct {
	Message aliyunChoiceMessage `json:"message"`
}

type aliyunResponse struct {
	Choices []aliyunChoice `json:"choices"`
}

func (p AliyunProvider) Generate(personaPrompt string, userText string) (string, error) {
	if p.Client == nil {
		// Robust transport with timeouts
		transport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 20 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
		}
		p.Client = &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		}
	}
	reqBody := aliyunRequest{
		Model: p.Model,
		Messages: []aliyunMessage{
			{Role: "system", Content: personaPrompt},
			{Role: "user", Content: userText},
		},
		MaxTokens:   512,
		Temperature: 0.7,
		TopP:        0.9,
	}
	b, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("%s/chat/completions", p.BaseURL)
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.APIKey)

	// Context timeout per request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	httpReq = httpReq.WithContext(ctx)

	// Retry on transient errors with exponential backoff
	resp, err := p.Client.Do(httpReq)
	if err != nil || resp.StatusCode >= 500 || resp.StatusCode == 429 {
		// Longer backoff for rate limiting
		backoffTime := 2 * time.Second
		if resp != nil && resp.StatusCode == 429 {
			backoffTime = 5 * time.Second
		}
		time.Sleep(backoffTime)
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		resp, err = p.Client.Do(httpReq)
		if err != nil {
			return "", err
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorMsg string
		switch resp.StatusCode {
		case 429:
			errorMsg = "rate limit exceeded (429)"
		case 401:
			errorMsg = "unauthorized (401) - check API key"
		case 402:
			errorMsg = "payment required (402) - insufficient credits"
		case 404:
			errorMsg = "model not found (404)"
		default:
			errorMsg = fmt.Sprintf("http error %d", resp.StatusCode)
		}
		return "", fmt.Errorf("aliyun %s", errorMsg)
	}
	var aliyunResp aliyunResponse
	if err := json.NewDecoder(resp.Body).Decode(&aliyunResp); err != nil {
		return "", err
	}
	if len(aliyunResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}
	out := aliyunResp.Choices[0].Message.Content
	// Short-answer mitigation: retry once with larger budget
	if len(out) < 20 {
		reqBody.MaxTokens = 768
		reqBody.Temperature = 0.6
		b2, _ := json.Marshal(reqBody)
		httpReq2, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(b2))
		httpReq2.Header = httpReq.Header.Clone()
		ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel2()
		httpReq2 = httpReq2.WithContext(ctx2)
		if resp2, err2 := p.Client.Do(httpReq2); err2 == nil {
			defer resp2.Body.Close()
			var aliyunResp2 aliyunResponse
			if err := json.NewDecoder(resp2.Body).Decode(&aliyunResp2); err == nil && len(aliyunResp2.Choices) > 0 {
				if c := aliyunResp2.Choices[0].Message.Content; len(c) >= len(out) {
					out = c
				}
			}
		}
	}
	return out, nil
}
