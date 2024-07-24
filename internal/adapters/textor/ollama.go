package textor

import "github.com/WildEgor/pi-storyteller/internal/configs"

var _ Textor = (*OllamaAdapter)(nil)

// OllamaAdapter wrapper around Ollama REST API
type OllamaAdapter struct {
	config *configs.OllamaConfig
}

// NewOllamaAdapter creates adapter
func NewOllamaAdapter(config *configs.OllamaConfig) *OllamaAdapter {
	return &OllamaAdapter{
		config,
	}
}

// Txt2Txt implements Textor.
func (o *OllamaAdapter) Txt2Txt(prompt string, opts *Opts) (result chan string, err error) {
	ch := make(chan string, 1)
	defer close(ch)

	return ch, nil
}
