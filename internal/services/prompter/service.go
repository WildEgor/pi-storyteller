package prompter

import (
	"encoding/json"
	"github.com/pemistahl/lingua-go"
	"os"
	"os/exec"
	"path/filepath"

	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"time"

	"github.com/WildEgor/pi-storyteller/internal/adapters/textor"
)

const parserPath = "scripts/parser.sh"

// Prompter ...
type Prompter struct {
	cache    *Cache
	detector lingua.LanguageDetector
	textor   textor.Textor
}

// New ...
func New(textor textor.Textor) *Prompter {
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

	return &Prompter{
		cache:    NewCache(),
		detector: detector,
		textor:   textor,
	}
}

// GetRandomNews ...
func (p *Prompter) GetRandomNews() (string, error) {
	wd, _ := os.Getwd()

	cmd := exec.Command("/bin/sh", filepath.Join(wd, parserPath))
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

	slog.Debug("", parsedNews)

	return parsedNews.Text, nil
}

// Random ...
func (p *Prompter) Random(lang string) []Conv {
	if len(lang) == 0 {
		lang = "en"
	}
	prompts := p.cache.Prompts(lang)
	//nolint
	randPrompt := prompts[rand.Intn(len(prompts))]
	story := p.cache.GetPrompt(randPrompt, lang)

	// TODO: find alternative for chat gpt
	// story, err := p.textor.Txt2Txt(prompt, nil)
	//if err != nil {
	//	return nil
	//}

	var conv []Conv
	for _, s := range strings.Split(story, ".") {
		if len(s) == 0 {
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

// WithPredefinedRandomStyle ...
func (p *Prompter) WithPredefinedRandomStyle(source string) []Conv {
	code, ok := p.detector.DetectLanguageOf(source)
	if !ok {
		slog.Warn("could not detect language", slog.Any("input", source))
		return nil
	}

	lang := strings.ToLower(code.IsoCode639_1().String())

	styles := p.cache.Styles(lang)
	//nolint
	randStyle := styles[rand.Intn(len(styles))]

	prompt := p.cache.GetStyle(randStyle, lang)

	var prompts []Conv
	for _, s := range strings.Split(source, ".") {
		if len(s) == 0 {
			continue
		}
		prompts = append(prompts, Conv{
			Style:    randStyle,
			Original: s,
			Prompt:   fmt.Sprintf(prompt, s),
		})
	}

	return prompts
}
