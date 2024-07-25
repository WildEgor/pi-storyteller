// Package imaginator provides abstraction around text2img network
package imaginator

import (
	"context"
	"image"
)

// Imagininator ...
type Imagininator interface {
	// GenerateImages receive prompt(s) and generate sequence of images
	GenerateImages(ctx context.Context, prompt []string) chan GeneratedImageResult
}

// GeneratedImageResult ...
type GeneratedImageResult struct {
	ID     int
	Image  image.Image
	Prompt string
	Error  error
}
