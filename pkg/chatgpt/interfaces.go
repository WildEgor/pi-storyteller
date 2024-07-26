package chatgpt

// Client for ChatGPT API
//
//go:generate mockery --name=Client --structname=ChatGPTClientMock --case=underscore
type Client interface {
	Generate(content string, opts *GenerateOpts) (string, error)
}

// GenerateOpts ...
type GenerateOpts struct {
	Model     string
	MaxTokens uint
}

// ClientConfig client config
type ClientConfig struct {
	BaseURL string
	Key     string
	Debug   bool
}

// ConfigFactory helper
type ConfigFactory func() *ClientConfig
