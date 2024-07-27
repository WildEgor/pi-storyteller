package templater

import (
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"
)

// templatePath ...
const templatePath = "assets/templates"

// Cache ...
type Cache struct {
	templates map[string]*template.Template
}

// NewCache ...
func NewCache() *Cache {
	c := &Cache{
		templates: make(map[string]*template.Template),
	}

	c.Init()

	return c
}

// Init ...
func (c *Cache) Init() {
	//nolint
	pwd, _ := os.Getwd()
	tp := filepath.Join(pwd, templatePath)

	files, err := os.ReadDir(tp)
	if err != nil {
		slog.Error("read templates error", slog.Any("err", err))
		panic(err)
	}

	c.templates = make(map[string]*template.Template, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		tml, err := template.ParseFiles(filepath.Join(tp, file.Name()))
		if err != nil {
			slog.Error("parse template error", slog.Any("err", err))
			continue
		}

		c.templates[file.Name()] = tml
	}
}

// Get ...
func (c *Cache) Get(name string) *template.Template {
	fn := fmt.Sprintf("%s.html", name)
	return c.templates[fn]
}
