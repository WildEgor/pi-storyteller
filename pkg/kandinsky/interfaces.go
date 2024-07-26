package kandinsky

import "context"

// Client for Kandinsky API
//
//go:generate mockery --name=Client --structname=KandinskyClientMock --case=underscore
type Client interface {
	GenerateImage(ctx context.Context, prompt string, opts *GenerateImageOpts) (*GenerateResult, error)
	GetModels(ctx context.Context) ([]ModelResult, error)
	GetTextToImageModel(ctx context.Context) (*ModelResult, error)
	CheckStatus(ctx context.Context, uuid string) (*GenerateImageResult, error)
}

// ModelType currently only supports one model type
type ModelType string

var (
	// TextToImage ...
	TextToImage ModelType = "TEXT2IMAGE"
)

// GenerateImageOpts settrings
type GenerateImageOpts struct {
	ModelId uint
	Width   uint
	Height  uint
	Neg     string
}

// ModelResult ...
type ModelResult struct {
	Id   uint      `json:"id"`
	Name string    `json:"name"`
	Type ModelType `json:"type"`
}

// GeneratePrompt ...
type GeneratePrompt struct {
	Prompt string `json:"query"`
}

// Type ...
type Type string

var (
	// Generate ...
	Generate Type = "GENERATE"
)

// Status ...
type Status string

var (
	// Initial when request image
	Initial Status = "INITIAL"
	// Done after generated
	Done Status = "DONE"
)

// GenerateParams ...
type GenerateParams struct {
	Type           Type           `json:"type"`
	NumImages      uint           `json:"numImages"`
	Width          uint           `json:"width"`
	Height         uint           `json:"height"`
	Neg            string         `json:"negativePromptUnclip"`
	GenerateParams GeneratePrompt `json:"generateParams"`
}

// GenerateData ...
type GenerateData struct {
	ModelId string         `json:"model_id"`
	Params  GenerateParams `json:"params"`
}

// GenerateResult ...
type GenerateResult struct {
	Status Status `json:"status"`
	UUID   string `json:"uuid"`
}

// GenerateImageResult ...
type GenerateImageResult struct {
	UUID     string   `json:"uuid"`
	Status   Status   `json:"status"`
	Images   []string `json:"images"`
	Censored bool     `json:"censored"`
}

// Done ...
func (r *GenerateImageResult) Done() bool {
	return r.Status == "DONE"
}

// ClientConfig client config
type ClientConfig struct {
	BaseURL string
	Key     string
	Secret  string
	Debug   bool
}

// ConfigFactory helper
type ConfigFactory func() *ClientConfig
