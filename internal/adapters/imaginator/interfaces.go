package imaginator

import "context"

type Imagininator interface {
	// GenerateImages recieve prompt(s) and generate sequence of images
	GenerateImages(ctx context.Context, prompt []string, result chan ImageGenerationResult, onUpdate func())
}
