package audio

import "testing"

func TestMockDualEngineCreation(t *testing.T) {
	engine := NewMockDualEngine()

	// Test initial state
	if engine.IsLoaded('A') {
		t.Error("Mock Deck A should not be loaded initially")
	}
	if engine.IsLoaded('B') {
		t.Error("Mock Deck B should not be loaded initially")
	}
	if engine.IsPlaying('A') {
		t.Error("Mock Deck A should not be playing initially")
	}
	if engine.IsPlaying('B') {
		t.Error("Mock Deck B should not be playing initially")
	}
}

func TestMockDualEngineWorkflow(t *testing.T) {
	engine := NewMockDualEngine()

	// Test loading
	err := engine.LoadToDeck("test.mp3", 'A')
	if err != nil {
		t.Fatalf("Failed to load to mock deck A: %v", err)
	}

	if !engine.IsLoaded('A') {
		t.Error("Mock deck A should be loaded")
	}

	if engine.IsPlaying('A') {
		t.Error("Mock deck A should not be playing after load")
	}

	// Test playing
	err = engine.PlayDeck('A')
	if err != nil {
		t.Fatalf("Failed to play mock deck A: %v", err)
	}

	if !engine.IsPlaying('A') {
		t.Error("Mock deck A should be playing")
	}

	// Test pausing
	err = engine.PauseDeck('A')
	if err != nil {
		t.Fatalf("Failed to pause mock deck A: %v", err)
	}

	if engine.IsPlaying('A') {
		t.Error("Mock deck A should not be playing after pause")
	}

	// Test that deck B is independent
	if engine.IsLoaded('B') {
		t.Error("Mock deck B should not be affected by deck A operations")
	}
}

func TestMockDualEngineInvalidOperations(t *testing.T) {
	engine := NewMockDualEngine()

	// Test playing without loading
	err := engine.PlayDeck('A')
	if err == nil {
		t.Error("Expected error when playing unloaded mock deck")
	}

	// Test invalid deck
	err = engine.LoadToDeck("test.mp3", 'X')
	if err == nil {
		t.Error("Expected error for invalid deck")
	}
}

func TestMockDualEngineVolumeControl(t *testing.T) {
	engine := NewMockDualEngine()

	// Test default volume
	if engine.GetVolume('A') != 1.0 {
		t.Errorf("Expected default volume 1.0, got %f", engine.GetVolume('A'))
	}

	// Test setting volume
	engine.SetVolume('A', 0.5)
	if engine.GetVolume('A') != 0.5 {
		t.Errorf("Expected volume 0.5, got %f", engine.GetVolume('A'))
	}

	// Test volume bounds
	engine.SetVolume('A', -1.0)
	if engine.GetVolume('A') != 0.0 {
		t.Errorf("Expected volume clamped to 0.0, got %f", engine.GetVolume('A'))
	}

	engine.SetVolume('A', 2.0)
	if engine.GetVolume('A') != 1.0 {
		t.Errorf("Expected volume clamped to 1.0, got %f", engine.GetVolume('A'))
	}
}

func TestMockDualEngineIndependentDecks(t *testing.T) {
	engine := NewMockDualEngine()

	// Load both decks
	engine.LoadToDeck("track1.mp3", 'A')
	engine.LoadToDeck("track2.mp3", 'B')

	// Play deck A
	engine.PlayDeck('A')

	// Verify independence
	if !engine.IsPlaying('A') {
		t.Error("Deck A should be playing")
	}
	if engine.IsPlaying('B') {
		t.Error("Deck B should not be playing")
	}

	// Set different volumes
	engine.SetVolume('A', 0.8)
	engine.SetVolume('B', 0.3)

	if engine.GetVolume('A') != 0.8 {
		t.Error("Deck A volume should be independent")
	}
	if engine.GetVolume('B') != 0.3 {
		t.Error("Deck B volume should be independent")
	}
}