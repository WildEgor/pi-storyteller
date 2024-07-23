package textor

type Template int

const (
	KindStory Template = iota
	FunnyStory
	AnimeStyleStory
)

// Opts adjustments
type Opts struct {
	SentencesLimit uint
	Template       Template
}

var DefaultOpts = &Opts{}