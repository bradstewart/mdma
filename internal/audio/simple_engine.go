package audio

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

type SimpleEngine struct {
	context *oto.Context
	player  oto.Player
	volume  float64
}

func NewSimpleEngine() (*SimpleEngine, error) {
	ctx, ready, err := oto.NewContext(44100, 2, 2)
	if err != nil {
		return nil, fmt.Errorf("failed to create audio context: %v", err)
	}
	<-ready

	return &SimpleEngine{
		context: ctx,
		volume:  1.0,
	}, nil
}

func (e *SimpleEngine) LoadAndPlay(filename string) error {
	if e.player != nil {
		e.player.Close()
	}

	ext := strings.ToLower(filepath.Ext(filename))

	var reader io.Reader

	switch ext {
	case ".mp3":
		file, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("failed to open MP3 file: %v", err)
		}

		decoder, err := mp3.NewDecoder(file)
		if err != nil {
			file.Close()
			return fmt.Errorf("failed to create MP3 decoder: %v", err)
		}
		reader = decoder

	case ".m4a":
		// For now, we'll need to use ffmpeg to convert M4A to raw PCM
		return e.loadM4AWithFFmpeg(filename)

	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	e.player = e.context.NewPlayer(reader)
	e.player.Play()

	return nil
}

func (e *SimpleEngine) loadM4AWithFFmpeg(filename string) error {
	// This is a simplified approach - in production we'd want proper streaming
	return fmt.Errorf("M4A support requires ffmpeg integration - not yet implemented")
}

func (e *SimpleEngine) Stop() {
	if e.player != nil {
		e.player.Pause()
	}
}

func (e *SimpleEngine) Play() {
	if e.player != nil {
		e.player.Play()
	}
}

func (e *SimpleEngine) IsPlaying() bool {
	return e.player != nil && e.player.IsPlaying()
}

func (e *SimpleEngine) Close() {
	if e.player != nil {
		e.player.Close()
	}
}