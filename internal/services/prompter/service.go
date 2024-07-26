package prompter

import (
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"time"

	"github.com/WildEgor/pi-storyteller/internal/configs"
	"github.com/pemistahl/lingua-go"
)

// Prompter ...
type Prompter struct {
	cache    *Cache
	detector lingua.LanguageDetector
}

// New ...
func New(appConfig *configs.AppConfig) *Prompter {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	languages := []lingua.Language{
		lingua.English,
		lingua.Russian,
	}

	detector := lingua.NewLanguageDetectorBuilder().
		FromLanguages(languages...).
		WithPreloadedLanguageModels().
		Build()

	return &Prompter{
		cache:    NewPromptsCache(appConfig.PromptsFilePath()),
		detector: detector,
	}
}

// Random ...
func (p *Prompter) Random(source string) []Conv {
	code, ok := p.detector.DetectLanguageOf(source)
	if !ok {
		slog.Warn("could not detect language", slog.Any("input", source))
		return nil
	}

	lang := strings.ToLower(code.IsoCode639_1().String())

	styles := p.cache.Keys(lang)
	randStyle := styles[rand.Intn(len(styles))]

	prompt := p.cache.Get(randStyle, lang)

	var prompts []Conv
	for _, s := range strings.Split(source, ".") {
		if len(s) == 0 {
			continue
		}
		prompts = append(prompts, Conv{
			Original: s,
			Prompt:   fmt.Sprintf(prompt, s),
		})
	}

	return prompts
}
