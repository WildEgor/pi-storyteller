package textor

import (
	"github.com/WildEgor/pi-storyteller/internal/configs"
	"github.com/WildEgor/pi-storyteller/pkg/chatgpt"
	"github.com/WildEgor/pi-storyteller/pkg/chatgpt/mocks"
)

// ChatGPTClientProvider wrapper for client
type ChatGPTClientProvider struct {
	chatgpt.Client
}

// NewChatGPTClientProvider create client
func NewChatGPTClientProvider(
	config *configs.ChatGPTConfig,
	appConfig *configs.AppConfig,
) *ChatGPTClientProvider {
	if appConfig.IsDebug() {
		return &ChatGPTClientProvider{mocks.NewChatGPTDummyClient(config.Config)}
	}
	return &ChatGPTClientProvider{chatgpt.New(config.Config)}
}
