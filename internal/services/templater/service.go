package templater

import (
	"bytes"
)

var _ Templater = (*Service)(nil)

// Service ...
type Service struct {
	cache *Cache
}

// New ...
func New() *Service {
	return &Service{
		cache: NewCache(),
	}
}

// Build ...
func (t *Service) Build(name string, data any) (string, error) {
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
