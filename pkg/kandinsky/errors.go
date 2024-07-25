package kandinsky

import "errors"

var (
	// ErrNoModels empty models
	ErrNoModels = errors.New("no models found")
	// ErrNoImage image not found
	ErrNoImage = errors.New("no images found")
)
