package imaginator

import "context"

type Imagininator interface {
	// GenerateImages recieve prompt(s) and generate sequence of images
	GenerateImages(ctx context.Context, prompt chan string, result chan ImageGenerationResult)
}
