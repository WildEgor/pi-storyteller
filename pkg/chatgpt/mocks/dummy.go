package mocks

import (
	"github.com/WildEgor/pi-storyteller/pkg/chatgpt"
)

// ChatGPTDummyClient ...
type ChatGPTDummyClient struct{}

// NewChatGPTDummyClient ...
func NewChatGPTDummyClient(config chatgpt.ConfigFactory) *ChatGPTDummyClient {
	return &ChatGPTDummyClient{}
}

// Generate ...
func (c *ChatGPTDummyClient) Generate(content string, opts *chatgpt.GenerateOpts) (string, error) {
	return "DUMMY TEST", nil
}
