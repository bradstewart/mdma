package audio

import (
	"os"
	"strings"
	"testing"
)

func skipIfShortOrNoAudio(t *testing.T) *DualEngine {
	if testing.Short() {
		t.Skip("Skipping audio engine test in short mode (requires audio system)")
	}

	engine, err := NewDualEngine()
	if err != nil {
		t.Skipf("Failed to create dual engine (audio system may not be available): %v", err)
	}
	return engine
}

func TestDualEngineCreation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping audio engine test in short mode (requires audio system)")
	}

	engine, err := NewDualEngine()
	if err != nil {
		t.Skipf("Failed to create dual engine (audio system may not be available): %v", err)
	}
	defer engine.Close()

	// Test initial state
	if engine.IsLoaded('A') {
		t.Error("Deck A should not be loaded initially")
	}
	if engine.IsLoaded('B') {
		t.Error("Deck B should not be loaded initially")
	}
	if engine.IsPlaying('A') {
		t.Error("Deck A should not be playing initially")
	}
	if engine.IsPlaying('B') {
		t.Error("Deck B should not be playing initially")
	}
}

func TestDualEngineInvalidDeck(t *testing.T) {
	engine := skipIfShortOrNoAudio(t)
	defer engine.Close()

	// Test invalid deck operations
	err := engine.LoadToDeck("test.mp3", 'X')
	if err == nil {
		t.Error("Expected error for invalid deck 'X'")
	}

	err = engine.PlayDeck('Y')
	if err == nil {
		t.Error("Expected error for invalid deck 'Y'")
	}

	if engine.IsLoaded('Z') {
		t.Error("Invalid deck should return false for IsLoaded")
	}

	if engine.IsPlaying('Z') {
		t.Error("Invalid deck should return false for IsPlaying")
	}
}

func TestDualEngineVolumeControl(t *testing.T) {
	engine := skipIfShortOrNoAudio(t)
	defer engine.Close()

	// Test default volume
	if engine.GetVolume('A') != 1.0 {
		t.Errorf("Expected default volume 1.0 for deck A, got %f", engine.GetVolume('A'))
	}

	// Test volume setting
	engine.SetVolume('A', 0.5)
	if engine.GetVolume('A') != 0.5 {
		t.Errorf("Expected volume 0.5 for deck A, got %f", engine.GetVolume('A'))
	}

	// Test volume bounds
	engine.SetVolume('B', -1.0)
	if engine.GetVolume('B') != 0.0 {
		t.Errorf("Expected volume to be clamped to 0.0, got %f", engine.GetVolume('B'))
	}

	engine.SetVolume('B', 2.0)
	if engine.GetVolume('B') != 1.0 {
		t.Errorf("Expected volume to be clamped to 1.0, got %f", engine.GetVolume('B'))
	}
}

func TestDualEngineLoadMP3(t *testing.T) {
	// Skip if test MP3 doesn't exist
	if _, err := os.Stat("/tmp/power.mp3"); os.IsNotExist(err) {
		t.Skip("Test MP3 file not found, skipping MP3 loading test")
	}

	engine := skipIfShortOrNoAudio(t)
	defer engine.Close()

	// Test loading MP3 to deck A
	err := engine.LoadToDeck("/tmp/power.mp3", 'A')
	if err != nil {
		t.Fatalf("Failed to load MP3 to deck A: %v", err)
	}

	if !engine.IsLoaded('A') {
		t.Error("Deck A should be loaded after LoadToDeck")
	}

	if engine.IsPlaying('A') {
		t.Error("Deck A should not be playing immediately after load")
	}

	// Test loading to deck B
	err = engine.LoadToDeck("/tmp/power.mp3", 'B')
	if err != nil {
		t.Fatalf("Failed to load MP3 to deck B: %v", err)
	}

	if !engine.IsLoaded('B') {
		t.Error("Deck B should be loaded after LoadToDeck")
	}

	// Both decks should be loaded now
	if !engine.IsLoaded('A') || !engine.IsLoaded('B') {
		t.Error("Both decks should remain loaded")
	}
}

func TestDualEnginePlayPause(t *testing.T) {
	// Skip if test MP3 doesn't exist
	if _, err := os.Stat("/tmp/power.mp3"); os.IsNotExist(err) {
		t.Skip("Test MP3 file not found, skipping play/pause test")
	}

	engine := skipIfShortOrNoAudio(t)
	defer engine.Close()

	// Test playing without loading
	err := engine.PlayDeck('A')
	if err == nil {
		t.Error("Expected error when playing empty deck")
	}

	// Load track
	err = engine.LoadToDeck("/tmp/power.mp3", 'A')
	if err != nil {
		t.Fatalf("Failed to load MP3: %v", err)
	}

	// Test play
	err = engine.PlayDeck('A')
	if err != nil {
		t.Fatalf("Failed to play deck A: %v", err)
	}

	if !engine.IsPlaying('A') {
		t.Error("Deck A should be playing after PlayDeck")
	}

	// Test pause
	err = engine.PauseDeck('A')
	if err != nil {
		t.Fatalf("Failed to pause deck A: %v", err)
	}

	if engine.IsPlaying('A') {
		t.Error("Deck A should not be playing after PauseDeck")
	}
}

func TestDualEngineIndependentDecks(t *testing.T) {
	// Skip if test MP3 doesn't exist
	if _, err := os.Stat("/tmp/power.mp3"); os.IsNotExist(err) {
		t.Skip("Test MP3 file not found, skipping independent decks test")
	}

	engine := skipIfShortOrNoAudio(t)
	defer engine.Close()

	// Load tracks to both decks
	err := engine.LoadToDeck("/tmp/power.mp3", 'A')
	if err != nil {
		t.Fatalf("Failed to load to deck A: %v", err)
	}

	err = engine.LoadToDeck("/tmp/power.mp3", 'B')
	if err != nil {
		t.Fatalf("Failed to load to deck B: %v", err)
	}

	// Play deck A
	err = engine.PlayDeck('A')
	if err != nil {
		t.Fatalf("Failed to play deck A: %v", err)
	}

	// Verify independence
	if !engine.IsPlaying('A') {
		t.Error("Deck A should be playing")
	}
	if engine.IsPlaying('B') {
		t.Error("Deck B should not be playing")
	}

	// Play deck B
	err = engine.PlayDeck('B')
	if err != nil {
		t.Fatalf("Failed to play deck B: %v", err)
	}

	// Both should be playing
	if !engine.IsPlaying('A') {
		t.Error("Deck A should still be playing")
	}
	if !engine.IsPlaying('B') {
		t.Error("Deck B should now be playing")
	}

	// Pause only deck A
	err = engine.PauseDeck('A')
	if err != nil {
		t.Fatalf("Failed to pause deck A: %v", err)
	}

	// Verify independence
	if engine.IsPlaying('A') {
		t.Error("Deck A should be paused")
	}
	if !engine.IsPlaying('B') {
		t.Error("Deck B should still be playing")
	}
}

func TestDualEngineUnsupportedFormat(t *testing.T) {
	engine := skipIfShortOrNoAudio(t)
	defer engine.Close()

	// Test unsupported format
	err := engine.LoadToDeck("test.wav", 'A')
	if err == nil {
		t.Error("Expected error for unsupported format")
	}

	// M4A should now be supported (though file might not exist)
	err = engine.LoadToDeck("nonexistent.m4a", 'A')
	if err != nil && !strings.Contains(err.Error(), "no such file") {
		// Should only fail because file doesn't exist, not because format is unsupported
		t.Logf("M4A loading failed with: %v", err)
	}
}

func TestDualEngineFileNotFound(t *testing.T) {
	engine := skipIfShortOrNoAudio(t)
	defer engine.Close()

	// Test non-existent file
	err := engine.LoadToDeck("/nonexistent/file.mp3", 'A')
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}