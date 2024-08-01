package prompter

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

// sourcePath default path
const sourcePath = "assets/source.json"

// Cache ...
type Cache struct {
	mu     sync.Mutex
	source *Source
}

// NewCache ...
func NewCache() *Cache {
	cache := &Cache{
		source: &Source{
			Actors: make([]Lang, 0),
			Places: make([]Lang, 0),
			Styles: make([]Lang, 0),
		},
	}

	cache.Init()

	return cache
}

// Init ...
func (t *Cache) Init() {
	t.mu.Lock()
	defer t.mu.Unlock()

	//nolint
	pwd, _ := os.Getwd()
	tp := filepath.Join(pwd, sourcePath)

	file, err := os.Open(tp)
	if err != nil {
		slog.Error("cannot open prompts", slog.Any("err", err))
		panic("")
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		slog.Error("cannot open prompts", slog.Any("err", err))
		panic("")
	}

	err = json.Unmarshal(byteValue, t.source)
	if err != nil {
		slog.Error("cannot open prompts", slog.Any("err", err))
		panic("")
	}
}

// Styles ...
func (t *Cache) Styles() []string {
	t.mu.Lock()
	defer t.mu.Unlock()

	keys := make([]string, 0, len(t.source.Styles))

	//nolint
	for _, v := range t.source.Styles {
		keys = append(keys, v.Alias)
	}

	return keys
}

// GetStyledPrompt ...
func (t *Cache) GetStyledPrompt(alias string, lang string, params ...string) string {
	t.mu.Lock()
	defer t.mu.Unlock()

	style, ok := lo.Find(t.source.Styles, func(item Lang) bool {
		return item.Alias == alias
	})
	if ok {
		if lang == "ru" {
			return fmt.Sprintf(style.Ru, params[0], params[1])
		}

		return fmt.Sprintf(style.En, params[0], params[1])
	}

	return ""
}

// Actors ...
func (t *Cache) Actors(lang string) []string {
	t.mu.Lock()
	defer t.mu.Unlock()

	keys := make([]string, 0, len(t.source.Actors))

	//nolint
	for _, v := range t.source.Actors {
		if lang == "ru" {
			keys = append(keys, v.Ru)
			continue
		}

		keys = append(keys, v.En)
	}

	return keys
}

// Places ...
func (t *Cache) Places(lang string) []string {
	t.mu.Lock()
	defer t.mu.Unlock()

	keys := make([]string, 0, len(t.source.Places))

	//nolint
	for _, v := range t.source.Places {
		if lang == "ru" {
			keys = append(keys, v.Ru)
			continue
		}

		keys = append(keys, v.En)
	}

	return keys
}
