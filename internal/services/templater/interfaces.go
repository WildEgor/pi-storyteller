package templater

// Templater ...
type Templater interface {
	Build(name string, data any) (string, error)
}

// SlideTemplate ...
type SlideTemplate struct {
	ImageBase64 string
	Description string
}
