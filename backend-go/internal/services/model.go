package services

import "fmt"

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
	// TODO: map roleID -> specific provider
	return MockProvider{}
}
