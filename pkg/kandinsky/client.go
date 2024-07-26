package kandinsky

import (
	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"

	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// client ...
type client struct {
	httpClient *resty.Client
}

// New creates http client for Kandinksy API
func New(config ConfigFactory) Client {
	httpClient := resty.New()

	httpClient.SetBaseURL(config().BaseURL).
		SetHeaders(map[string]string{
			"X-Key":    fmt.Sprintf("Key %s", config().Key),
			"X-Secret": fmt.Sprintf("Secret %s", config().Secret),
		}).SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		SetTimeout(5 * time.Second).
		SetDebug(config().Debug)

	return &client{
		httpClient,
	}
}

// GenerateImage request new image and get uuid
func (c *client) GenerateImage(ctx context.Context, prompt string, opts *GenerateImageOpts) (*GenerateResult, error) {
	params := GenerateParams{
		Type:      Generate,
		NumImages: 1,
		Width:     opts.Width,
		Height:    opts.Height,
		Neg:       opts.Neg,
		GenerateParams: GeneratePrompt{
			Prompt: prompt,
		},
	}

	paramsJson, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	//nolint
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "multipart/form-data").
		SetMultipartFormData(map[string]string{
			"model_id": strconv.Itoa(int(opts.ModelId)),
		}).
		SetMultipartField("params", "params.json", "application/json", bytes.NewReader(paramsJson)).
		Post("/key/api/v1/text2image/run")

	var generateResp *GenerateResult
	err = json.Unmarshal(resp.Body(), &generateResp)
	if err != nil {
		return nil, err
	}

	return generateResp, nil
}

// GetTextToImageModel helper method to find model
func (c *client) GetTextToImageModel(ctx context.Context) (*ModelResult, error) {
	models, err := c.GetModels(ctx)
	if err != nil {
		return nil, err
	}

	exitedModel, ok := lo.Find(models, func(m ModelResult) bool {
		return m.Type == TextToImage
	})
	if !ok {
		return nil, ErrNoModels
	}

	return &exitedModel, nil
}

// GetModels returns all available models
func (c *client) GetModels(ctx context.Context) ([]ModelResult, error) {
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		Get("/key/api/v1/models")
	if err != nil {
		return nil, err
	}

	var modelsResp []ModelResult
	err = json.Unmarshal(resp.Body(), &modelsResp)
	if err != nil {
		return nil, err
	}

	if len(modelsResp) == 0 {
		return nil, ErrNoModels
	}

	return modelsResp, nil
}

// CheckStatus of recent requested image. Return (once) generated image or fail
func (c *client) CheckStatus(ctx context.Context, uuid string) (*GenerateImageResult, error) {
	resp, err := c.httpClient.R().
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("/key/api/v1/text2image/status/%s", uuid))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, ErrNoImage
	}

	var imgResp GenerateImageResult
	err = json.Unmarshal(resp.Body(), &imgResp)
	if err != nil {
		return nil, err
	}

	return &imgResp, nil
}
