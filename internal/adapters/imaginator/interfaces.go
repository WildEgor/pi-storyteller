// Package imaginator provides abstraction around text2img network
package imaginator

import "context"

// Imagininator ...
type Imagininator interface {
	// GenerateImages receive prompt(s) and generate sequence of images
	GenerateImages(ctx context.Context, prompt []string) chan GeneratedImageResult
}
