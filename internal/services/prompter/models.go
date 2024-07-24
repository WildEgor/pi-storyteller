package prompter

type Style string

var (
	Realistic     Style = "realistic"
	Abstract      Style = "abstract"
	Surreal       Style = "surreal"
	Cyberpunk     Style = "cyberpunk"
	Fantasy       Style = "fantasy"
	Impressionist Style = "impressionist"
	Minimalist    Style = "minimalist"
	Vintage       Style = "vintage"
	Steampunk     Style = "steampunk"
	Anime         Style = "anime"
	Gothic        Style = "gothic"
	Noir          Style = "noir"
	PopArt        Style = "pop_art"
	Watercolor    Style = "watercolor"
	PixelArt      Style = "pixel_art"
	SciFi         Style = "sci_fi"
)

type Conv struct {
	Original string
	Prompt   string
}
