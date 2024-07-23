package kandinsky

type KandinskyModelType string

var (
	TextToImage KandinskyModelType = "TEXT2IMAGE"
)

type GenerateImageOpts struct {
	ModelId uint
	Width   uint
	Height  uint
}

type ModelResult struct {
	Id   uint               `json:"id"`
	Name string             `json:"name"`
	Type KandinskyModelType `json:"type"`
}

type GeneratePrompt struct {
	Prompt string `json:"query"`
}

type Type string

var (
	Generate Type = "GENERATE"
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
	Status string `json:"status"`
	UUID   string `json:"uuid"`
}

type GenerateImageResult struct {
	UUID   string   `json:"uuid"`
	Status string   `json:"status"`
	Images []string `json:"images"`
	//Censored bool     `json:"censored"`
	//GenTime  uint     `json:"generationTime"`
}

func (r *GenerateImageResult) Done() bool {
	return r.Status == "DONE"
}
