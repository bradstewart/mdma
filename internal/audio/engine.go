package audio

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/gordonklaus/portaudio"
	"github.com/hajimehoshi/go-mp3"
)

type Engine struct {
	stream    *portaudio.Stream
	isPlaying bool
	volume    float32
}

type AudioFile struct {
	reader   io.ReadCloser
	decoder  io.Reader
	filename string
}

func NewEngine() (*Engine, error) {
	err := portaudio.Initialize()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PortAudio: %v", err)
	}

	return &Engine{
		volume: 1.0,
	}, nil
}

func (e *Engine) LoadMP3(filename string) (*AudioFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filename, err)
	}

	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to create MP3 decoder: %v", err)
	}

	return &AudioFile{
		reader:   file,
		decoder:  decoder,
		filename: filename,
	}, nil
}

func (e *Engine) Play(audioFile *AudioFile) error {
	if e.isPlaying {
		return fmt.Errorf("audio engine is already playing")
	}

	// Create a simple callback that reads from the decoder
	buffer := make([]int16, 8192)

	stream, err := portaudio.OpenDefaultStream(
		0,     // input channels
		2,     // output channels (stereo)
		44100, // sample rate
		len(buffer)/2, // frames per buffer
		func(out [][]float32) {
			// Read samples from decoder
			for i := range out[0] {
				var sample int16
				// This is a simplified implementation - real implementation would need proper buffering
				if err := binary.Read(audioFile.decoder, binary.LittleEndian, &sample); err != nil {
					// End of file or error - fill with silence
					out[0][i] = 0
					out[1][i] = 0
				} else {
					// Convert int16 to float32 and apply volume
					floatSample := float32(sample) / 32768.0 * e.volume
					out[0][i] = floatSample
					out[1][i] = floatSample
				}
			}
		})

	if err != nil {
		return fmt.Errorf("failed to open audio stream: %v", err)
	}

	err = stream.Start()
	if err != nil {
		return fmt.Errorf("failed to start audio stream: %v", err)
	}

	e.stream = stream
	e.isPlaying = true
	return nil
}

func (e *Engine) Stop() error {
	if !e.isPlaying || e.stream == nil {
		return nil
	}

	err := e.stream.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop audio stream: %v", err)
	}

	err = e.stream.Close()
	if err != nil {
		return fmt.Errorf("failed to close audio stream: %v", err)
	}

	e.stream = nil
	e.isPlaying = false
	return nil
}

func (e *Engine) SetVolume(volume float32) {
	if volume < 0 {
		volume = 0
	} else if volume > 1 {
		volume = 1
	}
	e.volume = volume
}

func (e *Engine) IsPlaying() bool {
	return e.isPlaying
}

func (e *Engine) Close() error {
	if e.isPlaying {
		e.Stop()
	}
	return portaudio.Terminate()
}

func (af *AudioFile) Close() error {
	if af.reader != nil {
		return af.reader.Close()
	}
	return nil
}