package prompter

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Prompter struct {
	cache *Cache
}

func New() *Prompter {
	rand.Seed(time.Now().UnixNano())

	return &Prompter{
		cache: NewPromptsCache(),
	}
}

func (p *Prompter) Random(source string) []Conv {
	styles := p.cache.Keys()
	randStyle := styles[rand.Intn(len(styles))]
	prompt := p.cache.Get(randStyle)

	var prompts []Conv
	for _, s := range strings.Split(source, ".") {
		if len(s) == 0 {
			continue
		}
		prompts = append(prompts, Conv{
			Original: s,
			Prompt:   fmt.Sprintf(prompt, s),
		})
	}

	return prompts
}
