package imaginator

import "image"

type ImageGenerationResult struct {
	Prompt string
	Image  image.Image
}
