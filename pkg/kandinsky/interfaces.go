package kandinsky

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
