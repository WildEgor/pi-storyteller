package kandinsky

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

//go:generate mockery --name=IKandinskyClient --structname=KandinskyClientMock --case=underscore
type IKandinskyClient interface {
	GenerateImage(ctx context.Context, prompt string, opts *GenerateImageOpts) (*GenerateResult, error)
	GetModels(ctx context.Context) ([]ModelResult, error)
	CheckStatus(ctx context.Context, uuid string) (*GenerateImageResult, error)
}

type Client struct {
	client *resty.Client
}

func New(config IConfigFactory) *Client {
	client := resty.New()

	client.SetBaseURL(config().BaseURL).
		SetHeaders(map[string]string{
			"X-Key":    fmt.Sprintf("Key %s", config().Key),
			"X-Secret": fmt.Sprintf("Secret %s", config().Secret),
		}).SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second).
		SetTimeout(10 * time.Second)

	return &Client{
		client,
	}
}

func (c *Client) GenerateImage(ctx context.Context, prompt string, opts *GenerateImageOpts) (*GenerateResult, error) {
	params := GenerateParams{
		Type:      Generate,
		NumImages: 1,
		Width:     opts.Width,
		Height:    opts.Height,
		GenerateParams: GeneratePrompt{
			Prompt: prompt,
		},
	}

	paramsJson, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().
		SetContext(ctx).
		EnableTrace().
		SetDebug(true).
		SetHeader("Content-Type", "multipart/form-data").
		SetMultipartFormData(map[string]string{
			"model_id": strconv.Itoa(int(opts.ModelId)),
		}).
		SetMultipartField("params", "params.json", "application/json", bytes.NewReader(paramsJson)).
		SetDebug(true).
		Post("/key/api/v1/text2image/run")

	var generateResp *GenerateResult
	err = json.Unmarshal(resp.Body(), &generateResp)
	if err != nil {
		return nil, err
	}

	return generateResp, nil
}

func (c *Client) GetModels(ctx context.Context) ([]ModelResult, error) {
	resp, err := c.client.R().
		SetContext(ctx).
		EnableTrace().
		SetDebug(true).
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

func (c *Client) CheckStatus(ctx context.Context, uuid string) (*GenerateImageResult, error) {
	resp, err := c.client.R().
		EnableTrace().
		SetDebug(true).
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