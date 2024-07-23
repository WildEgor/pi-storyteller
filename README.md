General idea use open-source Ollama model that can be used for story generation and free models for images (for example, kandinsky).

Text generation must be generated using Raspberry Pi Zero 2 (need benchmarking) and using Go for image generation calls.
Also, Go server must render output html story files and send back to frontend.
Frontend simply show viewer (image corousel) with pagination and input field (if we want send custom promt) and "Generate" button (with cancel)

Golang server components:
- Imagininator: adapter for image generation
```go
type Imagininator interface {
    // GenerateImages recieve prompt(s) and generate sequence of images
    GenerateImages(prompt chan string, result chan ImageResult) error
}
```
- StoryGetter: adapter for story generation
```go
package story_adapter

type Template int

const (
    KindStory Template = iota
    FunnyStory
    AnimeStyleStory
)

// Opts adjustments
type Opts struct {
    SentencesLimit uint
    Template Template
}

// Adapter - for text2text generator (Ollama model)
type Adapter interface {
    // Txt2Txt send request to Ollama server and split (or not) to result chan or error
    Txt2Txt(promt string, opts *StoryOpts)(result chan string, err error)
}
```

StoryService as service for story generation
```go
type StorySlide struct {
    Page    int     `json:"page"`
    Content string  `json:"content"`
}

type StoryResult struct {
    Id string `json:"id"`
    Slides []StorySlide `json:"slides"`
    Status string `json:"status"`
}

type StoryStream struct {
    Id string
    slide chan StorySlide
}

type StoryService interface {
    Add(promt string) (sid string, err error)
    CheckStatus(sid string) (result *StoryResult, err error)
}
```

Handle Telegram Bot commands and JSON RPC and (soon)
Telegram Bot commands
```json
/start

/generate [template] [prompt] (where `prompt` your what about your story param and `template` predefined prompts)
```

- Handle one Task per User request (for telegram use nickname for json rpc api header);
- Set limit for Telegram Bot (/generate) about 5 Story per 30 minutes (I guess);
 
References:
- https://github.com/tvldz/storybook - using Python scripts and several models for displaying and ink display;
- https://github.com/vitoplantamura/OnnxStream - lightweight wrapper for text generation;
- github.com/go-telegram-bot-api/telegram-bot-api - Telegram Client fo Go;
