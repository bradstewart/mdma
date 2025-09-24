package audio

import (
	"testing"
)

func TestEngineCreation(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatalf("Failed to create audio engine: %v", err)
	}
	defer engine.Close()

	if engine.volume != 1.0 {
		t.Errorf("Expected default volume 1.0, got %f", engine.volume)
	}

	if engine.IsPlaying() {
		t.Error("Expected engine to not be playing initially")
	}
}

func TestVolumeControl(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatalf("Failed to create audio engine: %v", err)
	}
	defer engine.Close()

	// Test setting volume
	engine.SetVolume(0.5)
	if engine.volume != 0.5 {
		t.Errorf("Expected volume 0.5, got %f", engine.volume)
	}

	// Test volume bounds
	engine.SetVolume(-1.0)
	if engine.volume != 0.0 {
		t.Errorf("Expected volume to be clamped to 0.0, got %f", engine.volume)
	}

	engine.SetVolume(2.0)
	if engine.volume != 1.0 {
		t.Errorf("Expected volume to be clamped to 1.0, got %f", engine.volume)
	}
}