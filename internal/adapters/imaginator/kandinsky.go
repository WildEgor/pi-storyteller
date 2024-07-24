package imaginator

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"log/slog"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"

	"github.com/WildEgor/pi-storyteller/internal/configs"
	"github.com/WildEgor/pi-storyteller/pkg/kandinsky"
	"github.com/WildEgor/pi-storyteller/pkg/kandinsky/mocks"
)

var _ Imagininator = (*KandinskyAdapter)(nil)

// KandinskyClientProvider wrapper for client
type KandinskyClientProvider struct {
	kandinsky.Client
}

// NewKandinskyClientProvider create client
func NewKandinskyClientProvider(
	config *configs.KandinskyConfig,
) *KandinskyClientProvider {
	return &KandinskyClientProvider{kandinsky.New(config.Config)}
}

// NewKandinskyDummyClientProvider creates dummy client
func NewKandinskyDummyClientProvider(
	config *configs.KandinskyConfig,
) *KandinskyClientProvider {
	return &KandinskyClientProvider{mocks.NewKandinskyDummyClient(config.Config)}
}

// KandinskyAdapter around Kandinsky REST API
type KandinskyAdapter struct {
	client kandinsky.Client
	cache  *expirable.LRU[string, *kandinsky.ModelResult]
}

// NewKandinskyAdapter create adapter
func NewKandinskyAdapter(
	provider *KandinskyClientProvider,
) *KandinskyAdapter {
	cache := expirable.NewLRU[string, *kandinsky.ModelResult](5, nil, time.Hour*1)
	return &KandinskyAdapter{
		client: provider.Client,
		cache:  cache,
	}
}

// GenerateImages concurrency generate images by prompts
func (k *KandinskyAdapter) GenerateImages(ctx context.Context, prompts []string) chan GeneratedImageResult {
	ch := make(chan GeneratedImageResult, len(prompts))
	defer close(ch)

	slog.Info(fmt.Sprintf("prepared to generate n=%d images", len(prompts)))
	startAt := time.Now()

	var wg sync.WaitGroup
	for i, p := range prompts {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			uuid, err := k.generateImage(ctx, p)
			if err != nil {
				return
			}

			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			var errCounter uint8

			for {
				select {
				default:
				case <-ctx.Done():
					slog.Warn("context canceled")
					return
				case <-ticker.C:
					slog.Debug(fmt.Sprintf("process uuid %s", uuid))

					imgResult, cErr := k.client.CheckStatus(ctx, uuid)
					if cErr != nil {
						errCounter++
						if errCounter <= 3 {
							slog.Error("check status fail", slog.Any("err", err))
							continue
						}
						return
					}

					if !imgResult.Done() {
						continue
					}

					slog.Debug(fmt.Sprintf("done processing image=%d uuid=%s", id, uuid))

					img, _ := k.base64ToImage(imgResult.Images[0])

					ch <- GeneratedImageResult{
						ID:     id,
						Prompt: p,
						Image:  img,
					}
					return
				}
			}
		}(i)
	}
	wg.Wait()

	slog.Info(fmt.Sprintf("process took %s", time.Now().Sub(startAt)))

	return ch
}

// generateImage make request to Kandinsky for generation
func (k *KandinskyAdapter) generateImage(ctx context.Context, prompt string) (uuid string, err error) {
	existedModel, ok := k.cache.Get(string(kandinsky.TextToImage))
	if !ok {
		models, err := k.client.GetModels(ctx)
		if err != nil {
			return "", err
		}

		for _, model := range models {
			if model.Type == kandinsky.TextToImage {
				existedModel = &model
				break
			}
		}
	}
	if existedModel == nil {
		return "", ErrNoModels
	}

	k.cache.Add(string(kandinsky.TextToImage), existedModel)

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

// base64ToImage converts base64 string to golang image
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
