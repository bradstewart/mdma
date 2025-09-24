package audio

import "fmt"

// MockDualEngine provides a testable implementation without requiring audio hardware
type MockDualEngine struct {
	deckA *MockDeck
	deckB *MockDeck
}

type MockDeck struct {
	filename string
	isLoaded bool
	isPlaying bool
	volume   float64
}

func NewMockDualEngine() *MockDualEngine {
	return &MockDualEngine{
		deckA: &MockDeck{volume: 1.0},
		deckB: &MockDeck{volume: 1.0},
	}
}

func (e *MockDualEngine) LoadToDeck(filename string, deck rune) error {
	targetDeck := e.getMockDeck(deck)
	if targetDeck == nil {
		return fmt.Errorf("invalid deck: %c", deck)
	}

	// Simulate loading
	targetDeck.filename = filename
	targetDeck.isLoaded = true
	targetDeck.isPlaying = false

	return nil
}

func (e *MockDualEngine) PlayDeck(deck rune) error {
	targetDeck := e.getMockDeck(deck)
	if targetDeck == nil {
		return fmt.Errorf("invalid deck: %c", deck)
	}

	if !targetDeck.isLoaded {
		return fmt.Errorf("no track loaded on deck %c", deck)
	}

	targetDeck.isPlaying = true
	return nil
}

func (e *MockDualEngine) PauseDeck(deck rune) error {
	targetDeck := e.getMockDeck(deck)
	if targetDeck == nil {
		return fmt.Errorf("invalid deck: %c", deck)
	}

	if !targetDeck.isLoaded {
		return fmt.Errorf("no track loaded on deck %c", deck)
	}

	targetDeck.isPlaying = false
	return nil
}

func (e *MockDualEngine) StopDeck(deck rune) error {
	return e.PauseDeck(deck)
}

func (e *MockDualEngine) IsPlaying(deck rune) bool {
	targetDeck := e.getMockDeck(deck)
	if targetDeck == nil {
		return false
	}
	return targetDeck.isPlaying
}

func (e *MockDualEngine) IsLoaded(deck rune) bool {
	targetDeck := e.getMockDeck(deck)
	if targetDeck == nil {
		return false
	}
	return targetDeck.isLoaded
}

func (e *MockDualEngine) GetFilename(deck rune) string {
	targetDeck := e.getMockDeck(deck)
	if targetDeck == nil {
		return ""
	}
	return targetDeck.filename
}

func (e *MockDualEngine) SetVolume(deck rune, volume float64) {
	targetDeck := e.getMockDeck(deck)
	if targetDeck == nil {
		return
	}

	if volume < 0 {
		volume = 0
	} else if volume > 1 {
		volume = 1
	}
	targetDeck.volume = volume
}

func (e *MockDualEngine) GetVolume(deck rune) float64 {
	targetDeck := e.getMockDeck(deck)
	if targetDeck == nil {
		return 0
	}
	return targetDeck.volume
}

func (e *MockDualEngine) getMockDeck(deck rune) *MockDeck {
	switch deck {
	case 'A':
		return e.deckA
	case 'B':
		return e.deckB
	default:
		return nil
	}
}

func (e *MockDualEngine) Close() {
	// Nothing to clean up in mock
}