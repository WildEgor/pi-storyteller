package kandinsky

import "errors"

var (
	ErrNoModels = errors.New("no models found")
	ErrNoImage  = errors.New("no images found")
)
