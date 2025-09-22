package services

import "fmt"

// GenerateReply is a placeholder that simulates persona-aware responses.
func GenerateReply(personaPrompt string, userText string) string {
    if personaPrompt == "" {
        return fmt.Sprintf("(mock) You said: %s", userText)
    }
    return fmt.Sprintf("(mock persona) %s | Replying to: %s", personaPrompt, userText)
}


