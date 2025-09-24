package database

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDatabaseCreation(t *testing.T) {
	// Create temporary directory for test database
	tempDir, err := os.MkdirTemp("", "mdma-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override database path for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	db, err := NewDB()
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Verify database file was created
	dbPath := filepath.Join(tempDir, ".config", "mdma", "library.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}
}

func TestTrackOperations(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Test adding a track
	track := &Track{
		Path:        "/test/music/track1.mp3",
		Title:       "Test Track",
		Artist:      "Test Artist",
		Album:       "Test Album",
		Genre:       "Electronic",
		Year:        2023,
		Duration:    180.5,
		BPM:         128.0,
		Key:         "Am",
		EnergyLevel: 7,
		Rating:      4,
		CuePoints:   "[]",
		CustomTags:  "{}",
		Notes:       "Test notes",
	}

	err := db.AddTrack(track)
	if err != nil {
		t.Fatalf("Failed to add track: %v", err)
	}

	if track.ID == 0 {
		t.Error("Track ID should be set after adding to database")
	}

	// Test retrieving track by path
	retrieved, err := db.GetTrackByPath("/test/music/track1.mp3")
	if err != nil {
		t.Fatalf("Failed to get track by path: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Track should be found")
	}

	if retrieved.Title != "Test Track" {
		t.Errorf("Expected title 'Test Track', got '%s'", retrieved.Title)
	}

	if retrieved.BPM != 128.0 {
		t.Errorf("Expected BPM 128.0, got %f", retrieved.BPM)
	}

	// Test non-existent track
	notFound, err := db.GetTrackByPath("/nonexistent/track.mp3")
	if err != nil {
		t.Fatalf("Unexpected error for non-existent track: %v", err)
	}

	if notFound != nil {
		t.Error("Non-existent track should return nil")
	}
}

func TestTrackUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Add initial track
	track := &Track{
		Path:   "/test/update.mp3",
		Title:  "Original Title",
		Artist: "Original Artist",
		BPM:    120.0,
	}

	err := db.AddTrack(track)
	if err != nil {
		t.Fatalf("Failed to add track: %v", err)
	}

	// Update track (same path should replace)
	updatedTrack := &Track{
		Path:   "/test/update.mp3",
		Title:  "Updated Title",
		Artist: "Updated Artist",
		BPM:    140.0,
	}

	err = db.AddTrack(updatedTrack)
	if err != nil {
		t.Fatalf("Failed to update track: %v", err)
	}

	// Retrieve and verify update
	retrieved, err := db.GetTrackByPath("/test/update.mp3")
	if err != nil {
		t.Fatalf("Failed to get updated track: %v", err)
	}

	if retrieved.Title != "Updated Title" {
		t.Errorf("Expected updated title 'Updated Title', got '%s'", retrieved.Title)
	}

	if retrieved.BPM != 140.0 {
		t.Errorf("Expected updated BPM 140.0, got %f", retrieved.BPM)
	}
}

func TestGetAllTracks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Add multiple tracks
	tracks := []*Track{
		{Path: "/test/track1.mp3", Title: "Track 1", Artist: "Artist A", Album: "Album 1"},
		{Path: "/test/track2.mp3", Title: "Track 2", Artist: "Artist B", Album: "Album 1"},
		{Path: "/test/track3.mp3", Title: "Track 3", Artist: "Artist A", Album: "Album 2"},
	}

	for _, track := range tracks {
		err := db.AddTrack(track)
		if err != nil {
			t.Fatalf("Failed to add track: %v", err)
		}
	}

	// Get all tracks
	allTracks, err := db.GetAllTracks()
	if err != nil {
		t.Fatalf("Failed to get all tracks: %v", err)
	}

	if len(allTracks) != 3 {
		t.Errorf("Expected 3 tracks, got %d", len(allTracks))
	}

	// Verify sorting (should be by artist, album, title)
	if allTracks[0].Artist != "Artist A" {
		t.Errorf("Expected first track artist 'Artist A', got '%s'", allTracks[0].Artist)
	}
}

func TestSearchTracks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Add test tracks
	tracks := []*Track{
		{Path: "/test/electronic1.mp3", Title: "Synthesizer Dreams", Artist: "DJ Electronic", Album: "Digital Waves", Genre: "Electronic"},
		{Path: "/test/house1.mp3", Title: "House Party", Artist: "House Master", Album: "Club Nights", Genre: "House"},
		{Path: "/test/electronic2.mp3", Title: "Bass Drop", Artist: "DJ Electronic", Album: "Heavy Beats", Genre: "Electronic"},
		{Path: "/test/rock1.mp3", Title: "Guitar Solo", Artist: "Rock Band", Album: "Live Concert", Genre: "Rock"},
	}

	for _, track := range tracks {
		err := db.AddTrack(track)
		if err != nil {
			t.Fatalf("Failed to add track: %v", err)
		}
	}

	// Test search by artist
	results, err := db.SearchTracks("DJ Electronic")
	if err != nil {
		t.Fatalf("Failed to search tracks: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 tracks for 'DJ Electronic', got %d", len(results))
	}

	// Test search by genre
	results, err = db.SearchTracks("Electronic")
	if err != nil {
		t.Fatalf("Failed to search tracks: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 tracks for 'Electronic', got %d", len(results))
	}

	// Test search by title
	results, err = db.SearchTracks("House")
	if err != nil {
		t.Fatalf("Failed to search tracks: %v", err)
	}

	if len(results) != 1 { // Should match "House Party" title by "House Master" artist
		t.Errorf("Expected 1 track for 'House', got %d", len(results))
		for _, track := range results {
			t.Logf("Found track: %s - %s [%s] (%s)", track.Artist, track.Title, track.Album, track.Genre)
		}
	}

	// Test no results
	results, err = db.SearchTracks("Nonexistent")
	if err != nil {
		t.Fatalf("Failed to search tracks: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 tracks for 'Nonexistent', got %d", len(results))
	}
}

func TestTrackMetadata(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Test track with all metadata fields
	track := &Track{
		Path:        "/test/metadata.mp3",
		Title:       "Full Metadata Track",
		Artist:      "Metadata Artist",
		Album:       "Metadata Album",
		Genre:       "Test Genre",
		Year:        2023,
		Duration:    240.75,
		BPM:         132.5,
		Key:         "C#m",
		EnergyLevel: 8,
		Rating:      5,
		CuePoints:   `[{"name":"Intro","time":10.5},{"name":"Drop","time":60.0}]`,
		CustomTags:  `{"mood":"energetic","style":"progressive"}`,
		Notes:       "This track has a great breakdown at 2:30",
	}

	err := db.AddTrack(track)
	if err != nil {
		t.Fatalf("Failed to add track with metadata: %v", err)
	}

	// Retrieve and verify all fields
	retrieved, err := db.GetTrackByPath("/test/metadata.mp3")
	if err != nil {
		t.Fatalf("Failed to get track: %v", err)
	}

	if retrieved.Duration != 240.75 {
		t.Errorf("Expected duration 240.75, got %f", retrieved.Duration)
	}

	if retrieved.Key != "C#m" {
		t.Errorf("Expected key 'C#m', got '%s'", retrieved.Key)
	}

	if retrieved.EnergyLevel != 8 {
		t.Errorf("Expected energy level 8, got %d", retrieved.EnergyLevel)
	}

	if retrieved.CuePoints != `[{"name":"Intro","time":10.5},{"name":"Drop","time":60.0}]` {
		t.Errorf("Cue points not preserved correctly")
	}

	if retrieved.CustomTags != `{"mood":"energetic","style":"progressive"}` {
		t.Errorf("Custom tags not preserved correctly")
	}
}

func setupTestDB(t *testing.T) *DB {
	// Create temporary directory for test database
	tempDir, err := os.MkdirTemp("", "mdma-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Override database path for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)

	// Cleanup function to restore environment
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
		os.RemoveAll(tempDir)
	})

	db, err := NewDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	return db
}