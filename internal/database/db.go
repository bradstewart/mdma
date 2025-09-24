package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

type Track struct {
	ID          int64   `db:"id"`
	Path        string  `db:"path"`
	Title       string  `db:"title"`
	Artist      string  `db:"artist"`
	Album       string  `db:"album"`
	Genre       string  `db:"genre"`
	Year        int     `db:"year"`
	Duration    float64 `db:"duration"`
	BPM         float64 `db:"bpm"`
	Key         string  `db:"key"`
	EnergyLevel int     `db:"energy_level"`
	Rating      int     `db:"rating"`
	CuePoints   string  `db:"cue_points"` // JSON
	CustomTags  string  `db:"custom_tags"` // JSON
	Notes       string  `db:"notes"`
	LastPlayed  *string `db:"last_played"`
	PlayCount   int     `db:"play_count"`
	CreatedAt   string  `db:"created_at"`
	UpdatedAt   string  `db:"updated_at"`
}

type Playlist struct {
	ID        int64  `db:"id"`
	Name      string `db:"name"`
	Tracks    string `db:"tracks"` // JSON array of track IDs
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func NewDB() (*DB, error) {
	dbPath, err := getDBPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get database path: %v", err)
	}

	// Ensure database directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	db := &DB{conn: conn}

	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return db, nil
}

func (db *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS tracks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT UNIQUE NOT NULL,
		title TEXT,
		artist TEXT,
		album TEXT,
		genre TEXT,
		year INTEGER,
		duration REAL,
		bpm REAL,
		key TEXT,
		energy_level INTEGER DEFAULT 5,
		rating INTEGER DEFAULT 0,
		cue_points TEXT DEFAULT '[]',
		custom_tags TEXT DEFAULT '{}',
		notes TEXT DEFAULT '',
		last_played DATETIME,
		play_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS playlists (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		tracks TEXT DEFAULT '[]',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_tracks_artist ON tracks(artist);
	CREATE INDEX IF NOT EXISTS idx_tracks_album ON tracks(album);
	CREATE INDEX IF NOT EXISTS idx_tracks_genre ON tracks(genre);
	CREATE INDEX IF NOT EXISTS idx_tracks_bpm ON tracks(bpm);
	CREATE INDEX IF NOT EXISTS idx_tracks_key ON tracks(key);
	CREATE INDEX IF NOT EXISTS idx_tracks_energy_level ON tracks(energy_level);
	`

	_, err := db.conn.Exec(schema)
	return err
}

func (db *DB) AddTrack(track *Track) error {
	query := `
	INSERT OR REPLACE INTO tracks
	(path, title, artist, album, genre, year, duration, bpm, key, energy_level, rating, cue_points, custom_tags, notes, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`

	result, err := db.conn.Exec(query,
		track.Path, track.Title, track.Artist, track.Album, track.Genre,
		track.Year, track.Duration, track.BPM, track.Key, track.EnergyLevel,
		track.Rating, track.CuePoints, track.CustomTags, track.Notes)

	if err != nil {
		return err
	}

	track.ID, err = result.LastInsertId()
	return err
}

func (db *DB) GetTrackByPath(path string) (*Track, error) {
	query := `SELECT * FROM tracks WHERE path = ?`

	var track Track
	err := db.conn.QueryRow(query, path).Scan(
		&track.ID, &track.Path, &track.Title, &track.Artist, &track.Album,
		&track.Genre, &track.Year, &track.Duration, &track.BPM, &track.Key,
		&track.EnergyLevel, &track.Rating, &track.CuePoints, &track.CustomTags,
		&track.Notes, &track.LastPlayed, &track.PlayCount, &track.CreatedAt, &track.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &track, nil
}

func (db *DB) GetAllTracks() ([]*Track, error) {
	query := `SELECT * FROM tracks ORDER BY artist, album, title`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*Track
	for rows.Next() {
		var track Track
		err := rows.Scan(
			&track.ID, &track.Path, &track.Title, &track.Artist, &track.Album,
			&track.Genre, &track.Year, &track.Duration, &track.BPM, &track.Key,
			&track.EnergyLevel, &track.Rating, &track.CuePoints, &track.CustomTags,
			&track.Notes, &track.LastPlayed, &track.PlayCount, &track.CreatedAt, &track.UpdatedAt)

		if err != nil {
			return nil, err
		}
		tracks = append(tracks, &track)
	}

	return tracks, rows.Err()
}

func (db *DB) SearchTracks(query string) ([]*Track, error) {
	searchQuery := `
	SELECT * FROM tracks
	WHERE title LIKE ? OR artist LIKE ? OR album LIKE ? OR genre LIKE ?
	ORDER BY artist, album, title
	`

	searchTerm := "%" + query + "%"
	rows, err := db.conn.Query(searchQuery, searchTerm, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*Track
	for rows.Next() {
		var track Track
		err := rows.Scan(
			&track.ID, &track.Path, &track.Title, &track.Artist, &track.Album,
			&track.Genre, &track.Year, &track.Duration, &track.BPM, &track.Key,
			&track.EnergyLevel, &track.Rating, &track.CuePoints, &track.CustomTags,
			&track.Notes, &track.LastPlayed, &track.PlayCount, &track.CreatedAt, &track.UpdatedAt)

		if err != nil {
			return nil, err
		}
		tracks = append(tracks, &track)
	}

	return tracks, rows.Err()
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func getDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "mdma", "library.db"), nil
}