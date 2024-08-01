package prompter

// Prompter ...
type Prompter interface {
	GetRandomNews() (string, error)
	GetRandomStory(lang string) []Conv
	GetPredefinedRandomStyleStory(source string, sep bool) []Conv
}

// ParsedNews ...
type ParsedNews struct {
	Source string `json:"source"`
	Text   string `json:"text"`
	Link   string `json:"link"`
}

// Conv ...
type Conv struct {
	Style    string
	Original string
	Prompt   string
}

// Lang ...
type Lang struct {
	Alias string `json:"alias"`
	En    string `json:"en"`
	Ru    string `json:"ru"`
}

// Source ...
type Source struct {
	Actors []Lang `json:"actors"`
	Places []Lang `json:"places"`
	Styles []Lang `json:"styles"`
}
