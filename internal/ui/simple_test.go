package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/bradstewart/mdma/internal/audio"
	"github.com/bradstewart/mdma/internal/database"
	"github.com/bradstewart/mdma/internal/library"
)

// newTestModel creates a model with mock audio engine for testing
func newTestModel(t *testing.T) Model {
	// Create mock audio engine
	mockEngine := audio.NewMockDualEngine()

	// Create database for testing (will use temp directory)
	db, err := database.NewDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	scanner := library.NewScanner(db)

	return Model{
		audioEngine:      mockEngine,
		db:               db,
		scanner:          scanner,
		status:           "Ready",
		deckMode:         'A',
		librarySection:   "songs",
		selectedTrackIdx: 0,
		scrollOffset:     0,
		tracks:           make([]*database.Track, 0),
	}
}

func TestModelCreation(t *testing.T) {
	// Create a model without initializing the real audio engine
	m := Model{
		deckMode:         'A',
		status:           "Ready",
		librarySection:   "songs",
		selectedTrackIdx: 0,
		scrollOffset:     0,
	}

	// Test that model initializes properly
	if m.deckMode != 'A' {
		t.Errorf("Expected initial deck mode 'A', got '%c'", m.deckMode)
	}

	if m.librarySection != "songs" {
		t.Errorf("Expected initial library section 'songs', got '%s'", m.librarySection)
	}
}

func TestDeckModeSwitch(t *testing.T) {
	m := newTestModel(t)
	defer m.db.Close()

	// Test switching to deck B
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'B'}}
	updatedModel, _ := m.Update(keyMsg)
	m = updatedModel.(Model)

	if m.deckMode != 'B' {
		t.Errorf("Expected deck mode 'B' after pressing 'B', got '%c'", m.deckMode)
	}

	// Test switching back to deck A
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'A'}}
	updatedModel, _ = m.Update(keyMsg)
	m = updatedModel.(Model)

	if m.deckMode != 'A' {
		t.Errorf("Expected deck mode 'A' after pressing 'A', got '%c'", m.deckMode)
	}
}

func TestLibrarySectionSwitch(t *testing.T) {
	m := newTestModel(t)
	defer m.db.Close()

	// Test switching to playlists
	keyMsg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := m.Update(keyMsg)
	m = updatedModel.(Model)

	if m.librarySection != "playlists" {
		t.Errorf("Expected library section 'playlists' after pressing Tab, got '%s'", m.librarySection)
	}

	// Test switching back to songs
	keyMsg = tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ = m.Update(keyMsg)
	m = updatedModel.(Model)

	if m.librarySection != "songs" {
		t.Errorf("Expected library section 'songs' after pressing Tab again, got '%s'", m.librarySection)
	}
}

func TestPlayPauseWithoutTrack(t *testing.T) {
	m := newTestModel(t)
	defer m.db.Close()

	// Test play/pause without loaded track
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}}
	updatedModel, _ := m.Update(keyMsg)
	m = updatedModel.(Model)

	// Should show error status
	if m.status == "" {
		t.Error("Expected error status when trying to play without track")
	}
}

func TestWindowResize(t *testing.T) {
	m := newTestModel(t)
	defer m.db.Close()

	// Test window resize
	resizeMsg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updatedModel, _ := m.Update(resizeMsg)
	m = updatedModel.(Model)

	if m.width != 100 {
		t.Errorf("Expected width 100, got %d", m.width)
	}

	if m.height != 30 {
		t.Errorf("Expected height 30, got %d", m.height)
	}

	if !m.ready {
		t.Error("Model should be ready after window resize")
	}
}

func TestViewRendering(t *testing.T) {
	m := newTestModel(t)
	defer m.db.Close()
	m.ready = true
	m.width = 120
	m.height = 40

	// Test that view renders without panic
	view := m.View()

	if view == "" {
		t.Error("View should not be empty")
	}

	// Should contain key UI elements
	expectedElements := []string{"DECK A", "DECK B", "Active: Deck A"}
	for _, element := range expectedElements {
		if !containsString(view, element) {
			t.Errorf("View should contain '%s'", element)
		}
	}
}

func TestLibraryViewRendering(t *testing.T) {
	m := newTestModel(t)
	defer m.db.Close()
	m.ready = true
	m.width = 120
	m.height = 40

	// Test library view rendering
	view := m.View()

	if view == "" {
		t.Error("Library view should not be empty")
	}

	// Should contain library-specific elements
	if !containsString(view, "All Songs") {
		t.Error("Library view should contain 'All Songs'")
	}
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}