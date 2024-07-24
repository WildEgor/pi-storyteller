package kandinsky

type ClientConfig struct {
	BaseURL string
	Key     string
	Secret  string
	Debug   bool
}

type IConfigFactory func() *ClientConfig
