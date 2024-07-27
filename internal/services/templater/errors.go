package templater

import "errors"

var (
	// ErrTemplateNotFound ...
	ErrTemplateNotFound = errors.New("cannot find template")
	// ErrParseTemplate ...
	ErrParseTemplate = errors.New("cannot parse template")
)
