package kandinsky

type ClientConfig struct {
	BaseURL string
	Key     string
	Secret  string
}

type IConfigFactory func() *ClientConfig
