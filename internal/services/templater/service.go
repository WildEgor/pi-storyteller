package templater

import (
	"bytes"
	"errors"

	"github.com/WildEgor/pi-storyteller/internal/configs"
)

var (
	// ErrTemplateNotFound ...
	ErrTemplateNotFound = errors.New("cannot find template")
	// ErrParseTemplate ...
	ErrParseTemplate = errors.New("cannot parse template")
)

// Templater ...
type Templater struct {
	cache *TemplateCache
}

// NewTemplateService ...
func NewTemplateService(appConfig *configs.AppConfig) *Templater {
	cache := &TemplateCache{}
	cache.Init(appConfig.TemplatesPath())

	return &Templater{
		cache,
	}
}

// Build ...
func (t *Templater) Build(name string, data any) (string, error) {
	tmpl := t.cache.Get(name)
	if tmpl == nil {
		return "", ErrTemplateNotFound
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		return "", ErrParseTemplate
	}

	return buf.String(), nil
}
