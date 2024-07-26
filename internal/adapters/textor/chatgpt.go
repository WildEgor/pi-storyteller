package textor

import (
	"github.com/WildEgor/pi-storyteller/internal/configs"
	"github.com/WildEgor/pi-storyteller/pkg/chatgpt"
	"github.com/WildEgor/pi-storyteller/pkg/chatgpt/mocks"
)

var _ Textor = (*ChatGPTAdapter)(nil)

// ChatGPTClientProvider wrapper for client
type ChatGPTClientProvider struct {
	chatgpt.Client
}

// NewChatGPTClientProvider create client
func NewChatGPTClientProvider(
	config *configs.ChatGPTConfig,
) *ChatGPTClientProvider {
	return &ChatGPTClientProvider{chatgpt.New(config.Config)}
}

// NewChatGPTDummyClientProvider creates dummy client
func NewChatGPTDummyClientProvider(
	config *configs.ChatGPTConfig,
) *ChatGPTClientProvider {
	return &ChatGPTClientProvider{mocks.NewChatGPTDummyClient(config.Config)}
}

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
