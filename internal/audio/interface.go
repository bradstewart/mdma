package audio

// AudioEngine defines the interface for dual-deck audio operations
type AudioEngine interface {
	LoadToDeck(filename string, deck rune) error
	PlayDeck(deck rune) error
	PauseDeck(deck rune) error
	StopDeck(deck rune) error
	IsPlaying(deck rune) bool
	IsLoaded(deck rune) bool
	GetFilename(deck rune) string
	SetVolume(deck rune, volume float64)
	GetVolume(deck rune) float64
	Close()
}