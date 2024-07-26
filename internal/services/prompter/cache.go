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

// stylesPath default path
const stylesPath = "assets/styles.json"

// promptsPath default path
const promptsPath = "assets/prompts.json"

// Cache ...
type Cache struct {
	mu      sync.Mutex
	styles  map[string]map[string]string
	prompts map[string]map[string]string
}

// NewCache ...
func NewCache() *Cache {
	cache := &Cache{
		styles:  make(map[string]map[string]string),
		prompts: make(map[string]map[string]string),
	}

	cache.styles["en"] = make(map[string]string)
	cache.styles["ru"] = make(map[string]string)
	cache.prompts["en"] = make(map[string]string)
	cache.prompts["ru"] = make(map[string]string)

	cache.Init()

	return cache
}

// Init ...
func (t *Cache) Init() {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i, p := range []string{stylesPath, promptsPath} {
		data := t.readAndMapJSON(p)

		for lang, v := range data {
			if strings.Contains(lang, "_en") {
				if i == 0 {
					t.styles["en"][strings.ReplaceAll(lang, "_en", "")] = v
				} else {
					t.prompts["en"][strings.ReplaceAll(lang, "_en", "")] = v
				}
			}
			if strings.Contains(lang, "_ru") {
				if i == 0 {
					t.styles["ru"][strings.ReplaceAll(lang, "_ru", "")] = v
				} else {
					t.prompts["ru"][strings.ReplaceAll(lang, "_ru", "")] = v
				}
			}
		}
	}
}

func (t *Cache) readAndMapJSON(path string) map[string]string {
	//nolint
	pwd, _ := os.Getwd()
	tp := filepath.Join(pwd, path)

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

	return prompts
}

// GetStyle ...
func (t *Cache) GetStyle(name string, lang string) string {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(lang) == 0 {
		lang = "en"
	}

	return t.styles[lang][name]
}

// Styles ...
func (t *Cache) Styles(lang string) (keys []string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(lang) == 0 {
		lang = "en"
	}

	//nolint
	for k, _ := range t.styles[lang] {
		keys = append(keys, k)
	}

	return keys
}

// GetPrompt ...
func (t *Cache) GetPrompt(name string, lang string) string {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(lang) == 0 {
		lang = "en"
	}

	return t.prompts[lang][name]
}

// Prompts ...
func (t *Cache) Prompts(lang string) (keys []string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(lang) == 0 {
		lang = "en"
	}

	//nolint
	for k, _ := range t.prompts[lang] {
		keys = append(keys, k)
	}

	return keys
}
