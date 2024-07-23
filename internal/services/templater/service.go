package templater

import (
	"bytes"
	"errors"
)

var (
	ErrTemplateNotFound = errors.New("cannot find template")
	ErrParseTemplate    = errors.New("cannot parse template")
)

type Templater struct {
	cache *TemplateCache
}

func NewTemplateService() *Templater {
	cache := &TemplateCache{}
	cache.Init()

	return &Templater{
		cache,
	}
}

func (t *Templater) Build(name string, data interface{}) (string, error) {
	tmpl := t.cache.GetByName(name)
	if tmpl == nil {
		return "", ErrTemplateNotFound
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		return "", ErrParseTemplate
	}

	return buf.String(), nil
}
