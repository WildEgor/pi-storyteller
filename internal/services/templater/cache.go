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

// TemplateCache ...
type TemplateCache struct {
	templates map[string]*template.Template
}

// Init ...
func (t *TemplateCache) Init(path string) {
	//nolint
	pwd, _ := os.Getwd()
	tp := filepath.Join(pwd, templatePath)

	if len(path) != 0 {
		tp = filepath.Join(path)
	}

	files, err := os.ReadDir(tp)
	if err != nil {
		slog.Error("read templates error", slog.Any("err", err))
		panic(err)
	}

	t.templates = make(map[string]*template.Template, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		tml, err := template.ParseFiles(filepath.Join(tp, file.Name()))
		if err != nil {
			slog.Error("parse template error", slog.Any("err", err))
			continue
		}

		t.templates[file.Name()] = tml
	}
}

// Get ...
func (t *TemplateCache) Get(name string) *template.Template {
	fn := fmt.Sprintf("%s.html", name)

	return t.templates[fn]
}
