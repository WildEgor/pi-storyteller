package imaginator

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"time"

	"github.com/WildEgor/pi-storyteller/internal/configs"
	"github.com/WildEgor/pi-storyteller/pkg/kandinsky"
)

var _ Imagininator = (*KandinskyAdapter)(nil)

type KandinskyAdapter struct {
	client *kandinsky.Client
}

func NewKandinskyAdapter(
	config *configs.KandinskyConfig,
) *KandinskyAdapter {
	return &KandinskyAdapter{
		client: kandinsky.New(config.Config),
	}
}

// GenerateImages implements Imagininator.
func (k *KandinskyAdapter) GenerateImages(ctx context.Context, prompt chan string, result chan ImageGenerationResult) {
	go func() {
		defer close(result)

		for p := range prompt {
			runner := func() {
				uuid, err := k.generateImage(ctx, p)
				if err != nil {
					return
				}
				ticker := time.NewTicker(5 * time.Second)

				for {
					select {
					case <-ctx.Done():
						ticker.Stop()
						return
					case <-ticker.C:
						imgResult, err := k.client.CheckStatus(ctx, uuid)
						if err != nil {
							// TODO: count retry and cancel
							continue
						}

						if imgResult.Done() {
							img, _ := k.base64ToImage(imgResult.Images[0])

							result <- ImageGenerationResult{
								Prompt: p,
								Image:  img,
							}

							ticker.Stop()
							return
						}
					}
				}
			}
			runner()
		}
	}()
}

func (k *KandinskyAdapter) generateImage(ctx context.Context, prompt string) (uuid string, err error) {
	models, err := k.client.GetModels(ctx)
	if err != nil {
		return "", err
	}

	var existedModel *kandinsky.ModelResult
	for _, model := range models {
		if model.Type == kandinsky.TextToImage {
			existedModel = &model
			break
		}
	}

	if existedModel == nil {
		return "", ErrNoModels
	}

	resp, err := k.client.GenerateImage(ctx, prompt, &kandinsky.GenerateImageOpts{
		ModelId: existedModel.Id,
		Width:   512,
		Height:  512,
	})
	if err != nil {
		return "", err
	}

	return resp.UUID, nil
}

func (k *KandinskyAdapter) base64ToImage(b64 string) (result image.Image, err error) {
	imageData, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}

	return img, nil
}
