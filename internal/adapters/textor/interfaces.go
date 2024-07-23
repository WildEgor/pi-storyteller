package textor

// ITextor - for text2text generator
type ITextor interface {
	// Txt2Txt send request
	Txt2Txt(promt string, opts *Opts) (result chan string, err error)
}
