package textor

// Textor for text2text generator
type Textor interface {
	// Txt2Txt send request and receive any text
	Txt2Txt(prompt string, opts *Opts) (result string, err error)
}

// Template WithPredefinedRandomStyle template
type Template int

const (
	// KindStory ...
	KindStory Template = iota
	// FunnyStory ...
	FunnyStory
	// AnimeStyleStory ...
	AnimeStyleStory
)

// Opts adjustments
type Opts struct {
	SentencesLimit uint
	Template       Template
}

// DefaultOpts ...
var DefaultOpts = &Opts{}
