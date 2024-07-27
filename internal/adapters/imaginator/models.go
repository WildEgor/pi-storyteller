package imaginator

import "errors"

var (
	// ErrNoModels ...
	ErrNoModels = errors.New("no models found")
	// ErrNoImage ...
	ErrNoImage = errors.New("no images found")
)
