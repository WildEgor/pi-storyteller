package prompter

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

const promptsPath = "assets/prompts.json"

// Cache ...
type Cache struct {
	mu    sync.Mutex
	cache map[string]string
}

func NewPromptsCache(path string) *Cache {
	cache := &Cache{}
	cache.Init(path)

	return cache
}

// Init ...
func (t *Cache) Init(path string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	pwd, _ := os.Getwd()
	tp := filepath.Join(pwd, promptsPath)
	if len(path) != 0 {
		tp = filepath.Join(path)
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

	// Unmarshal the JSON data into the map
	err = json.Unmarshal(byteValue, &t.cache)
	if err != nil {
		slog.Error("cannot open prompts", slog.Any("err", err))
		panic("")
	}
}

// Get ...
func (t *Cache) Get(name string) string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.cache[name]
}

// Keys ...
func (t *Cache) Keys() (keys []string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for k, _ := range t.cache {
		keys = append(keys, k)
	}

	return keys
}
