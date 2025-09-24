package library

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bradstewart/mdma/internal/database"
)

func TestScannerCreation(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	scanner := NewScanner(db)
	if scanner == nil {
		t.Fatal("Scanner should not be nil")
	}

	if scanner.db != db {
		t.Error("Scanner should reference the provided database")
	}
}

func TestIsAudioFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"test.mp3", true},
		{"test.MP3", true},
		{"test.m4a", true},
		{"test.flac", true},
		{"test.wav", true},
		{"test.ogg", true},
		{"test.aac", true},
		{"test.txt", false},
		{"test.jpg", false},
		{"test.pdf", false},
		{"test", false},
		{"", false},
	}

	for _, test := range tests {
		result := isAudioFile(test.filename)
		if result != test.expected {
			t.Errorf("isAudioFile(%s) = %v, expected %v", test.filename, result, test.expected)
		}
	}
}

func TestParseBPM(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"128", 128.0},
		{"128.5", 128.5},
		{"140.25", 140.25},
		{"invalid", 0},
		{"", 0},
		{"999", 0}, // Out of range
		{"30", 0},  // Too low
	}

	for _, test := range tests {
		result := parseBPM(test.input)
		if result != test.expected {
			t.Errorf("parseBPM(%s) = %f, expected %f", test.input, result, test.expected)
		}
	}
}

func TestParseEnergyLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"1", 1},
		{"5", 5},
		{"10", 10},
		{"0", 5},  // Out of range, should default
		{"11", 5}, // Out of range, should default
		{"invalid", 5},
		{"", 5},
	}

	for _, test := range tests {
		result := parseEnergyLevel(test.input)
		if result != test.expected {
			t.Errorf("parseEnergyLevel(%s) = %d, expected %d", test.input, result, test.expected)
		}
	}
}

func TestScanEmptyDirectory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create empty temp directory
	tempDir, err := os.MkdirTemp("", "mdma-scan-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	scanner := NewScanner(db)
	err = scanner.ScanDirectory(tempDir)
	if err != nil {
		t.Fatalf("Failed to scan empty directory: %v", err)
	}

	// Should have no tracks
	tracks, err := db.GetAllTracks()
	if err != nil {
		t.Fatalf("Failed to get tracks: %v", err)
	}

	if len(tracks) != 0 {
		t.Errorf("Expected 0 tracks in empty directory, got %d", len(tracks))
	}
}

func TestScanDirectoryWithMixedFiles(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create temp directory with mixed file types
	tempDir, err := os.MkdirTemp("", "mdma-scan-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files (empty files for testing)
	files := []string{
		"track1.mp3",
		"track2.flac",
		"document.txt",
		"image.jpg",
		"track3.ogg",
		"readme.md",
	}

	for _, filename := range files {
		file, err := os.Create(filepath.Join(tempDir, filename))
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
		file.Close()
	}

	scanner := NewScanner(db)
	err = scanner.ScanDirectory(tempDir)
	if err != nil {
		t.Fatalf("Failed to scan directory: %v", err)
	}

	// Should have found 3 audio files
	tracks, err := db.GetAllTracks()
	if err != nil {
		t.Fatalf("Failed to get tracks: %v", err)
	}

	if len(tracks) != 3 {
		t.Errorf("Expected 3 audio tracks, got %d", len(tracks))
	}

	// Verify the correct files were found
	expectedFiles := []string{"track1.mp3", "track2.flac", "track3.ogg"}
	foundFiles := make(map[string]bool)

	for _, track := range tracks {
		filename := filepath.Base(track.Path)
		foundFiles[filename] = true
	}

	for _, expected := range expectedFiles {
		if !foundFiles[expected] {
			t.Errorf("Expected to find audio file %s", expected)
		}
	}
}

func TestScanNonexistentDirectory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	scanner := NewScanner(db)
	err := scanner.ScanDirectory("/nonexistent/directory")
	if err == nil {
		t.Error("Expected error when scanning nonexistent directory")
	}
}

func TestScanSubdirectories(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create nested directory structure
	tempDir, err := os.MkdirTemp("", "mdma-scan-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create subdirectories
	artistDir := filepath.Join(tempDir, "Artist Name")
	albumDir := filepath.Join(artistDir, "Album Name")
	err = os.MkdirAll(albumDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectories: %v", err)
	}

	// Create test files in different levels
	files := []string{
		filepath.Join(tempDir, "root_track.mp3"),
		filepath.Join(artistDir, "artist_track.mp3"),
		filepath.Join(albumDir, "album_track.mp3"),
	}

	for _, filename := range files {
		file, err := os.Create(filename)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
		file.Close()
	}

	scanner := NewScanner(db)
	err = scanner.ScanDirectory(tempDir)
	if err != nil {
		t.Fatalf("Failed to scan directory: %v", err)
	}

	// Should have found all 3 files
	tracks, err := db.GetAllTracks()
	if err != nil {
		t.Fatalf("Failed to get tracks: %v", err)
	}

	if len(tracks) != 3 {
		t.Errorf("Expected 3 tracks from nested directories, got %d", len(tracks))
	}
}

func TestRescanUpdatedFile(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "mdma-scan-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file
	testFile := filepath.Join(tempDir, "test.mp3")
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	scanner := NewScanner(db)

	// First scan
	err = scanner.ScanDirectory(tempDir)
	if err != nil {
		t.Fatalf("Failed to scan directory: %v", err)
	}

	tracks, err := db.GetAllTracks()
	if err != nil {
		t.Fatalf("Failed to get tracks: %v", err)
	}

	if len(tracks) != 1 {
		t.Fatalf("Expected 1 track after first scan, got %d", len(tracks))
	}

	firstScanTime := tracks[0].UpdatedAt

	// Second scan (file unchanged) - should be skipped
	err = scanner.ScanDirectory(tempDir)
	if err != nil {
		t.Fatalf("Failed to rescan directory: %v", err)
	}

	tracks, err = db.GetAllTracks()
	if err != nil {
		t.Fatalf("Failed to get tracks after rescan: %v", err)
	}

	if len(tracks) != 1 {
		t.Fatalf("Expected 1 track after rescan, got %d", len(tracks))
	}

	// Update time should be the same (file wasn't re-processed)
	if tracks[0].UpdatedAt != firstScanTime {
		t.Error("Track should not be re-processed if file hasn't changed")
	}
}

func setupTestDB(t *testing.T) *database.DB {
	// Create temporary directory for test database
	tempDir, err := os.MkdirTemp("", "mdma-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Override database path for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)

	// Cleanup function
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
		os.RemoveAll(tempDir)
	})

	db, err := database.NewDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	return db
}