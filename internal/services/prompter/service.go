package prompter

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pemistahl/lingua-go"

	"github.com/WildEgor/pi-storyteller/internal/adapters/textor"
)

var _ Prompter = (*Service)(nil)

const parserPath = "scripts/parser.sh"

// Service ...
type Service struct {
	cache    *Cache
	detector lingua.LanguageDetector
	gpt      textor.Textor
}

// New ...
func New(gpt textor.Textor) *Service {
	//nolint
	rand.New(rand.NewSource(time.Now().UnixNano()))

	languages := []lingua.Language{
		lingua.English,
		lingua.Russian,
	}

	detector := lingua.NewLanguageDetectorBuilder().
		FromLanguages(languages...).
		WithPreloadedLanguageModels().
		Build()

	return &Service{
		cache:    NewCache(),
		detector: detector,
		gpt:      gpt,
	}
}

// GetRandomNews ...
func (p *Service) GetRandomNews() (string, error) {
	//nolint
	wd, _ := os.Getwd()

	//nolint
	cmd := exec.Command("/bin/bash", filepath.Join(wd, parserPath))
	if cmd.Err != nil {
		return "", cmd.Err
	}

	raw, err := cmd.Output()
	if err != nil {
		return "", err
	}

	parsedNews := &ParsedNews{}

	pErr := json.Unmarshal(raw, parsedNews)
	if pErr != nil {
		return "", pErr
	}

	return fmt.Sprintf("%s \n Link: %s", parsedNews.Text, parsedNews.Link), nil
}

// GetRandomStory ...
func (p *Service) GetRandomStory(lang string) []Conv {
	if lang != "ru" && lang != "en" {
		lang = "en"
	}

	prompts := p.cache.Prompts(lang)
	//nolint
	randPrompt := prompts[rand.Intn(len(prompts))]
	story := p.cache.GetPrompt(randPrompt, lang)

	var conv []Conv
	for _, s := range strings.Split(story, ".") {
		if len(s) <= 3 {
			continue
		}
		conv = append(conv, Conv{
			Style:    randPrompt,
			Original: s,
			Prompt:   s,
		})
	}

	return conv
}

// GetPredefinedRandomStyleStory ...
func (p *Service) GetPredefinedRandomStyleStory(source string, sep bool) []Conv {
	code, ok := p.detector.DetectLanguageOf(source)
	if !ok {
		slog.Warn("could not detect language", slog.Any("input", source))
		return nil
	}

	lang := strings.ToLower(code.IsoCode639_1().String())
	if lang != "ru" && lang != "en" {
		lang = "en"
	}

	styles := p.cache.Styles(lang)
	//nolint
	randStyle := styles[rand.Intn(len(styles))]

	prompt := p.cache.GetStyle(randStyle, lang)

	var prompts []Conv
	if sep {
		for _, s := range strings.Split(source, ".") {
			if len(s) <= 3 {
				continue
			}
			prompts = append(prompts, Conv{
				Style:    randStyle,
				Original: s,
				Prompt:   fmt.Sprintf(prompt, s),
			})
		}
	} else {
		prompts = append(prompts, Conv{
			Style:    randStyle,
			Original: source,
			Prompt:   fmt.Sprintf(prompt, source),
		})
	}

	return prompts
}
