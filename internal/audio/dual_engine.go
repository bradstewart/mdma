package audio

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

type DualEngine struct {
	context *oto.Context
	deckA   *Deck
	deckB   *Deck
}

type Deck struct {
	player   oto.Player
	volume   float64
	filename string
	isLoaded bool
}

func NewDualEngine() (*DualEngine, error) {
	ctx, ready, err := oto.NewContext(44100, 2, 2)
	if err != nil {
		return nil, fmt.Errorf("failed to create audio context: %v", err)
	}
	<-ready

	return &DualEngine{
		context: ctx,
		deckA: &Deck{
			volume: 1.0,
		},
		deckB: &Deck{
			volume: 1.0,
		},
	}, nil
}

func (e *DualEngine) LoadToDeck(filename string, deck rune) error {
	var targetDeck *Deck
	switch deck {
	case 'A':
		targetDeck = e.deckA
	case 'B':
		targetDeck = e.deckB
	default:
		return fmt.Errorf("invalid deck: %c", deck)
	}

	// Stop and close existing player
	if targetDeck.player != nil {
		targetDeck.player.Close()
		targetDeck.player = nil
	}

	ext := strings.ToLower(filepath.Ext(filename))
	var reader io.Reader

	switch ext {
	case ".mp3":
		file, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("failed to open MP3 file: %v", err)
		}
		// Note: Don't close file here - the decoder needs it
		defer func() {
			if err != nil {
				file.Close()
			}
		}()

		decoder, err := mp3.NewDecoder(file)
		if err != nil {
			return fmt.Errorf("failed to create MP3 decoder: %v", err)
		}
		reader = decoder

	case ".m4a", ".aac", ".mp4":
		// Use FFmpeg to convert M4A/AAC to PCM data
		cmd := exec.Command("ffmpeg",
			"-i", filename,
			"-f", "s16le",     // 16-bit signed little endian PCM
			"-ar", "44100",    // 44.1kHz sample rate
			"-ac", "2",        // stereo
			"-")               // output to stdout

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to create ffmpeg pipe: %v", err)
		}

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start ffmpeg: %v", err)
		}

		reader = stdout

	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	targetDeck.player = e.context.NewPlayer(reader)
	targetDeck.filename = filename
	targetDeck.isLoaded = true

	return nil
}

func (e *DualEngine) PlayDeck(deck rune) error {
	targetDeck := e.getDeck(deck)
	if targetDeck == nil {
		return fmt.Errorf("invalid deck: %c", deck)
	}

	if targetDeck.player == nil {
		return fmt.Errorf("no track loaded on deck %c", deck)
	}

	targetDeck.player.Play()
	return nil
}

func (e *DualEngine) PauseDeck(deck rune) error {
	targetDeck := e.getDeck(deck)
	if targetDeck == nil {
		return fmt.Errorf("invalid deck: %c", deck)
	}

	if targetDeck.player == nil {
		return fmt.Errorf("no track loaded on deck %c", deck)
	}

	targetDeck.player.Pause()
	return nil
}

func (e *DualEngine) StopDeck(deck rune) error {
	targetDeck := e.getDeck(deck)
	if targetDeck == nil {
		return fmt.Errorf("invalid deck: %c", deck)
	}

	if targetDeck.player == nil {
		return fmt.Errorf("no track loaded on deck %c", deck)
	}

	targetDeck.player.Pause()
	// Reset to beginning (oto/v2 doesn't have a direct reset, so we'll reload)
	return nil
}

func (e *DualEngine) IsPlaying(deck rune) bool {
	targetDeck := e.getDeck(deck)
	if targetDeck == nil || targetDeck.player == nil {
		return false
	}

	return targetDeck.player.IsPlaying()
}

func (e *DualEngine) IsLoaded(deck rune) bool {
	targetDeck := e.getDeck(deck)
	if targetDeck == nil {
		return false
	}

	return targetDeck.isLoaded
}

func (e *DualEngine) GetFilename(deck rune) string {
	targetDeck := e.getDeck(deck)
	if targetDeck == nil {
		return ""
	}

	return targetDeck.filename
}

func (e *DualEngine) SetVolume(deck rune, volume float64) {
	targetDeck := e.getDeck(deck)
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

func (e *DualEngine) GetVolume(deck rune) float64 {
	targetDeck := e.getDeck(deck)
	if targetDeck == nil {
		return 0
	}

	return targetDeck.volume
}

func (e *DualEngine) getDeck(deck rune) *Deck {
	switch deck {
	case 'A':
		return e.deckA
	case 'B':
		return e.deckB
	default:
		return nil
	}
}

func (e *DualEngine) Close() {
	if e.deckA.player != nil {
		e.deckA.player.Close()
	}
	if e.deckB.player != nil {
		e.deckB.player.Close()
	}
}