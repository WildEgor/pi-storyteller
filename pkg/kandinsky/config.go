package kandinsky

// ClientConfig client config
type ClientConfig struct {
	BaseURL string
	Key     string
	Secret  string
	Debug   bool
}

// IConfigFactory helper
type IConfigFactory func() *ClientConfig
