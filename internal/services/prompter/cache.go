package prompter

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// promptsPath default path
const promptsPath = "assets/prompts.json"

// Cache ...
type Cache struct {
	mu         sync.Mutex
	dictionary map[string]map[string]string
}

// NewPromptsCache ...
func NewPromptsCache(path string) *Cache {
	cache := &Cache{
		dictionary: make(map[string]map[string]string),
	}

	cache.dictionary["en"] = make(map[string]string)
	cache.dictionary["ru"] = make(map[string]string)

	cache.Init(path)

	return cache
}

// Init ...
func (t *Cache) Init(path string) {
	//nolint
	pwd, _ := os.Getwd()
	tp := filepath.Join(pwd, promptsPath)
	if len(path) != 0 {
		tp = path
	}

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

	prompts := make(map[string]string)

	// Unmarshal the JSON data into the map
	err = json.Unmarshal(byteValue, &prompts)
	if err != nil {
		slog.Error("cannot open prompts", slog.Any("err", err))
		panic("")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	for lang, v := range prompts {
		if strings.Contains(lang, "_en") {
			t.dictionary["en"][strings.ReplaceAll(lang, "_en", "")] = v
		}

		if strings.Contains(lang, "_ru") {
			t.dictionary["ru"][strings.ReplaceAll(lang, "_ru", "")] = v
		}
	}
}

// Get ...
func (t *Cache) Get(name string, lang string) string {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(lang) == 0 {
		lang = "en"
	}

	return t.dictionary[lang][name]
}

// Keys ...
func (t *Cache) Keys(lang string) (keys []string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(lang) == 0 {
		lang = "en"
	}

	//nolint
	for k, _ := range t.dictionary[lang] {
		keys = append(keys, k)
	}

	return keys
}
