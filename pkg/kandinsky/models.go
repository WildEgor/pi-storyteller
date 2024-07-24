package kandinsky

type ModelType string

var (
	TextToImage ModelType = "TEXT2IMAGE"
)

type GenerateImageOpts struct {
	ModelId uint
	Width   uint
	Height  uint
}

type ModelResult struct {
	Id   uint      `json:"id"`
	Name string    `json:"name"`
	Type ModelType `json:"type"`
}

type GeneratePrompt struct {
	Prompt string `json:"query"`
}

type Type string

var (
	Generate Type = "GENERATE"
)

type Status string

var (
	Initial Status = "INITIAL"
	Done    Status = "DONE"
)

type GenerateParams struct {
	Type           Type           `json:"type"`
	NumImages      uint           `json:"numImages"`
	Width          uint           `json:"width"`
	Height         uint           `json:"height"`
	GenerateParams GeneratePrompt `json:"generateParams"`
}

type GenerateData struct {
	ModelId string         `json:"model_id"`
	Params  GenerateParams `json:"params"`
}

type GenerateResult struct {
	Status Status `json:"status"`
	UUID   string `json:"uuid"`
}

type GenerateImageResult struct {
	UUID   string   `json:"uuid"`
	Status Status   `json:"status"`
	Images []string `json:"images"`
	//Censored bool     `json:"censored"`
	//GenTime  uint     `json:"generationTime"`
}

func (r *GenerateImageResult) Done() bool {
	return r.Status == "DONE"
}
