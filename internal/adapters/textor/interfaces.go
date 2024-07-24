package textor

// Textor for text2text generator
type Textor interface {
	// Txt2Txt send request and receive any text
	Txt2Txt(prompt string, opts *Opts) (result chan string, err error)
}
