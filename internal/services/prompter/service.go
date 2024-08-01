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
	actors := p.cache.Actors(lang)
	//nolint
	randActor := actors[rand.Intn(len(actors))]

	places := p.cache.Places(lang)
	//nolint
	randPlace := places[rand.Intn(len(places))]

	style := p.cache.Styles()
	randStyle := style[rand.Intn(len(style))]

	//nolint
	story := p.cache.GetStyledPrompt(randStyle, lang, randActor, randPlace)

	var conv []Conv
	for _, s := range strings.Split(story, ".") {
		if len(s) <= 3 {
			continue
		}
		conv = append(conv, Conv{
			Style:    randStyle,
			Original: s,
			Prompt:   s,
		})
	}

	return conv
}

// GetPredefinedRandomStyleStory ...
func (p *Service) GetPredefinedRandomStyleStory(source string, sep bool) []Conv {
	// TODO: remove heavy lang detection
	code, ok := p.detector.DetectLanguageOf(source)
	if !ok {
		slog.Warn("could not detect language", slog.Any("input", source))
		return nil
	}

	lang := strings.ToLower(code.IsoCode639_1().String())
	if lang != "ru" && lang != "en" {
		lang = "en"
	}

	styles := p.cache.Styles()
	//nolint
	randStyle := styles[rand.Intn(len(styles))]

	var prompts []Conv
	if sep {
		for _, s := range strings.Split(source, ".") {
			if len(s) <= 3 {
				continue
			}
			prompts = append(prompts, Conv{
				Style:    randStyle,
				Original: s,
				Prompt:   p.cache.GetStyledPrompt(randStyle, lang, s),
			})
		}
	} else {
		prompts = append(prompts, Conv{
			Style:    randStyle,
			Original: source,
			Prompt:   p.cache.GetStyledPrompt(randStyle, lang, source),
		})
	}

	return prompts
}
