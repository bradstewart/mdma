package library

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhowden/tag"
	"github.com/bradstewart/mdma/internal/database"
)

type Scanner struct {
	db *database.DB
}

func NewScanner(db *database.DB) *Scanner {
	return &Scanner{db: db}
}

func (s *Scanner) ScanDirectory(rootPath string) error {
	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !isAudioFile(path) {
			return nil
		}

		// Check if track already exists in database
		existingTrack, err := s.db.GetTrackByPath(path)
		if err != nil {
			return fmt.Errorf("error checking existing track: %v", err)
		}

		// Skip if file hasn't been modified since last scan
		if existingTrack != nil {
			modTime := info.ModTime().Format(time.RFC3339)
			if existingTrack.UpdatedAt >= modTime {
				return nil
			}
		}

		// Extract metadata and add/update track
		track, err := s.extractMetadata(path, info)
		if err != nil {
			fmt.Printf("Warning: Failed to extract metadata from %s: %v\n", path, err)
			// Still add the track with basic info
			track = &database.Track{
				Path:  path,
				Title: filepath.Base(path),
			}
		}

		err = s.db.AddTrack(track)
		if err != nil {
			return fmt.Errorf("error adding track to database: %v", err)
		}

		fmt.Printf("Indexed: %s - %s\n", track.Artist, track.Title)
		return nil
	})
}

func (s *Scanner) extractMetadata(filePath string, info os.FileInfo) (*database.Track, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	metadata, err := tag.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	track := &database.Track{
		Path:     filePath,
		Title:    metadata.Title(),
		Artist:   metadata.Artist(),
		Album:    metadata.Album(),
		Genre:    metadata.Genre(),
		Duration: 0, // Would need audio analysis for accurate duration
	}

	// Handle year
	if year := metadata.Year(); year != 0 {
		track.Year = year
	}

	// Handle track number
	if trackNum, _ := metadata.Track(); trackNum != 0 {
		// Could store track number if needed
	}

	// Set defaults if metadata is missing
	if track.Title == "" {
		// Use filename without extension as title
		name := filepath.Base(filePath)
		ext := filepath.Ext(name)
		track.Title = strings.TrimSuffix(name, ext)
	}

	if track.Artist == "" {
		track.Artist = "Unknown Artist"
	}

	if track.Album == "" {
		track.Album = "Unknown Album"
	}

	// Try to extract BPM and Key from custom tags (common in DJ software)
	for key, value := range metadata.Raw() {
		keyLower := strings.ToLower(fmt.Sprintf("%v", key))
		valueStr := fmt.Sprintf("%v", value)

		switch {
		case strings.Contains(keyLower, "bpm") || strings.Contains(keyLower, "tempo"):
			// Try to parse BPM
			if bpm := parseBPM(valueStr); bpm > 0 {
				track.BPM = bpm
			}
		case strings.Contains(keyLower, "key") || strings.Contains(keyLower, "initialkey"):
			track.Key = valueStr
		case strings.Contains(keyLower, "energy"):
			if energy := parseEnergyLevel(valueStr); energy > 0 {
				track.EnergyLevel = energy
			}
		}
	}

	return track, nil
}

func isAudioFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".mp3", ".m4a", ".flac", ".wav", ".ogg", ".aac":
		return true
	}
	return false
}

func parseBPM(value string) float64 {
	// Simple BPM parsing - in real implementation would be more robust
	var bpm float64
	if _, err := fmt.Sscanf(value, "%f", &bpm); err == nil {
		if bpm >= 60 && bpm <= 300 { // Reasonable BPM range
			return bpm
		}
	}
	return 0
}

func parseEnergyLevel(value string) int {
	// Parse energy level (1-10 scale)
	var energy int
	if _, err := fmt.Sscanf(value, "%d", &energy); err == nil {
		if energy >= 1 && energy <= 10 {
			return energy
		}
	}
	return 5 // Default energy level
}