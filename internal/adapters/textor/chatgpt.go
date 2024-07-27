package textor

import (
	"github.com/WildEgor/pi-storyteller/pkg/chatgpt"
)

var _ Textor = (*ChatGPTAdapter)(nil)

// ChatGPTAdapter wrapper around ChatGPT REST API
type ChatGPTAdapter struct {
	client chatgpt.Client
}

// NewChatGPTAdapter creates adapter
func NewChatGPTAdapter(provider *ChatGPTClientProvider) *ChatGPTAdapter {
	return &ChatGPTAdapter{
		client: provider.Client,
	}
}

// Txt2Txt ...
func (o *ChatGPTAdapter) Txt2Txt(prompt string, opts *Opts) (result string, err error) {
	res, err := o.client.Generate(prompt, nil)
	if err != nil {
		return "", err
	}

	return res, nil
}
