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

const (
	// MaxGenErrors Default error counter
	MaxGenErrors = 3
	// DefaultImgSize Generated img size
	DefaultImgSize = 512
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
				case <-ctx.Done():
					slog.Warn("context canceled")
					return
				case <-ticker.C:
					slog.Debug(fmt.Sprintf("processing uuid=%s", uuid))

					imgResult, cErr := k.client.CheckStatus(ctx, uuid)
					if cErr != nil {
						errCounter++
						if errCounter <= MaxGenErrors {
							slog.Error("check status fail", slog.Any("err", err))
							continue
						}

						ch <- GeneratedImageResult{
							ID:     id,
							Prompt: p,
							Error:  cErr,
						}
						return
					}

					if !imgResult.Done() {
						continue
					}

					slog.Debug(fmt.Sprintf("done processing image=%d uuid=%s", id, uuid))

					img, convErr := k.base64ToImage(imgResult.Images[0])

					ch <- GeneratedImageResult{
						ID:     id,
						Prompt: p,
						Image:  img,
						Error:  convErr,
					}
					return
				}
			}
		}(i)
	}
	wg.Wait()

	slog.Info(fmt.Sprintf("process done, took - %s", time.Since(startAt)))

	return ch
}

// generateImage make request to Kandinsky for generation
func (k *KandinskyAdapter) generateImage(ctx context.Context, prompt string) (uuid string, err error) {
	existedModel, ok := k.cache.Get(string(kandinsky.TextToImage))
	if !ok {
		existedModel, err = k.client.GetTextToImageModel(ctx)
		if err != nil {
			return "", err
		}
	}

	k.cache.Add(string(kandinsky.TextToImage), existedModel)

	resp, err := k.client.GenerateImage(ctx, prompt, &kandinsky.GenerateImageOpts{
		ModelId: existedModel.Id,
		Width:   DefaultImgSize,
		Height:  DefaultImgSize,
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
