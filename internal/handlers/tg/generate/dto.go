package tg_generate_handler

import "fmt"

// GenerateCommandDTO ...
type GenerateCommandDTO struct {
	Nickname  string
	ChatID    string
	MessageID string
	Prompt    string
}

// Validate ...
func (p GenerateCommandDTO) Validate() error {
	if len(p.Prompt) == 0 {
		return fmt.Errorf("empty prompt not allowed")
	}

	if len(p.Prompt) < 10 {
		return fmt.Errorf("too many short prompt")
	}

	if len(p.Prompt) > 2024 {
		return fmt.Errorf("too many long prompt")
	}

	return nil
}
