package chatgpt

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

var _ Client = (*client)(nil)

// client ...
type client struct {
	httpClient *resty.Client
}

// New creates http client for ChatGPT API
func New(config ConfigFactory) Client {
	httpClient := resty.New()

	httpClient.SetBaseURL(config().BaseURL).
		SetHeader("Content-Type", "application/json").
		SetAuthToken(config().Key).
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		SetTimeout(5 * time.Second).
		SetDebug(config().Debug)

	return &client{
		httpClient,
	}
}

// Generate ...
func (c *client) Generate(content string, opts *GenerateOpts) (string, error) {
	model := "babbage-002"
	mx := 50

	if opts != nil {
		if len(opts.Model) != 0 {
			model = opts.Model
		}
		if opts.MaxTokens != 0 {
			mx = int(opts.MaxTokens)
		}
	}

	response, err := c.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"model": model,
			"messages": []any{map[string]any{
				"role":    "system",
				"content": content,
			},
			},
			"max_tokens": mx,
		}).
		Post("/v1/chat/completions")

	if err != nil {
		log.Fatalf("Error while sending send the request: %v", err)
	}

	body := response.Body()

	var data map[string]any
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	// Extract the content from the JSON response
	result := data["choices"].([]any)[0].(map[string]any)["message"].(map[string]any)["content"].(string)

	return result, nil
}
