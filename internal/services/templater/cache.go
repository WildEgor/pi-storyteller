package templater

import (
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"path"
	"path/filepath"
)

const templatePath = "templates"

type TemplateCache struct {
	templates map[string]*template.Template
}

func (t *TemplateCache) Init() {
	pwd, _ := os.Getwd()
	tp := filepath.Join(pwd, templatePath)

	files, err := os.ReadDir(tp)
	if err != nil {
		slog.Error("read templates error", slog.Any("err", err))
		panic(err)
	}

	t.templates = make(map[string]*template.Template, len(files))
	for _, file := range files {
		tml, err := template.ParseFiles(path.Join(templatePath, file.Name()))
		if err != nil {
			slog.Error("parse template error", slog.Any("err", err))
			continue
		}

		t.templates[file.Name()] = tml
	}
}

func (t *TemplateCache) GetByName(name string) *template.Template {
	fn := fmt.Sprintf("%s.html", name)

	return t.templates[fn]
}
