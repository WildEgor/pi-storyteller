package mocks

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log/slog"
	"time"

	"github.com/WildEgor/pi-storyteller/pkg/kandinsky"
)

// KandinskyDummyClient ...
type KandinskyDummyClient struct{}

// NewKandinskyDummyClient ...
func NewKandinskyDummyClient(config kandinsky.ConfigFactory) *KandinskyDummyClient {
	return &KandinskyDummyClient{}
}

// GenerateImage ...
func (c *KandinskyDummyClient) GenerateImage(ctx context.Context, prompt string, opts *kandinsky.GenerateImageOpts) (*kandinsky.GenerateResult, error) {
	time.Sleep(5 * time.Second)

	return &kandinsky.GenerateResult{
		Status: kandinsky.Done,
		UUID:   "777",
	}, nil
}

// GetTextToImageModel ...
func (c *KandinskyDummyClient) GetTextToImageModel(ctx context.Context) (*kandinsky.ModelResult, error) {
	return &kandinsky.ModelResult{
		Id:   1,
		Type: kandinsky.TextToImage,
	}, nil
}

// GetModels ...
func (c *KandinskyDummyClient) GetModels(ctx context.Context) ([]kandinsky.ModelResult, error) {
	return []kandinsky.ModelResult{{
		Id:   1,
		Type: kandinsky.TextToImage,
	}}, nil
}

// CheckStatus ...
func (c *KandinskyDummyClient) CheckStatus(ctx context.Context, uuid string) (*kandinsky.GenerateImageResult, error) {
	width := 100
	height := 100
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	red := color.RGBA{R: 125, G: 125, B: 125, A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{C: red}, image.Point{}, draw.Src)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		slog.Error("failed to encode image", slog.Any("err", err))
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	time.Sleep(10 * time.Second)

	return &kandinsky.GenerateImageResult{
		UUID:   uuid,
		Status: kandinsky.Done,
		Images: []string{encoded},
	}, nil
}
