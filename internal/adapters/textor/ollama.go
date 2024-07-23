package textor

import "github.com/WildEgor/pi-storyteller/internal/configs"

var _ ITextor = (*OllamaAdapter)(nil)

type OllamaAdapter struct {
	config *configs.OllamaConfig
}

func NewOllamaAdapter(config *configs.OllamaConfig) *OllamaAdapter {
	return &OllamaAdapter{
		config,
	}
}

// Txt2Txt implements ITextor.
func (o *OllamaAdapter) Txt2Txt(promt string, opts *Opts) (result chan string, err error) {
	panic("unimplemented")
}
