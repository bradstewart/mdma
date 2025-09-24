package ui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/bradstewart/mdma/internal/audio"
	"github.com/bradstewart/mdma/internal/database"
	"github.com/bradstewart/mdma/internal/library"
)

type Model struct {
	width              int
	height             int
	ready              bool
	audioEngine        audio.AudioEngine
	db                 *database.DB
	scanner            *library.Scanner
	deckATrack         *database.Track
	deckBTrack         *database.Track
	status             string
	deckMode           rune // 'A' or 'B' - current active deck

	// Library browser state
	librarySection     string   // "songs" or "playlists"
	selectedTrackIdx   int
	scrollOffset       int
	tracks             []*database.Track
	searchQuery        string
	searching          bool
}

type tickMsg struct{}

type SearchMsg struct {
	query string
}

type TracksLoadedMsg struct {
	tracks []*database.Track
}

func NewModel() Model {
	engine, err := audio.NewDualEngine()
	if err != nil {
		return Model{status: fmt.Sprintf("Audio error: %v", err)}
	}

	db, err := database.NewDB()
	if err != nil {
		return Model{status: fmt.Sprintf("Database error: %v", err)}
	}

	scanner := library.NewScanner(db)

	// Load existing tracks from database
	tracks, err := db.GetAllTracks()
	if err != nil {
		tracks = make([]*database.Track, 0)
	}

	return Model{
		audioEngine:      engine,
		db:               db,
		scanner:          scanner,
		status:           "Ready",
		deckMode:         'A', // Start with deck A active
		librarySection:   "songs",
		selectedTrackIdx: 0,
		scrollOffset:     0,
		tracks:           tracks,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case TracksLoadedMsg:
		m.tracks = msg.tracks
		m.selectedTrackIdx = 0
		m.scrollOffset = 0
		m.status = fmt.Sprintf("Loaded %d tracks", len(msg.tracks))
		return m, nil

	case SearchMsg:
		return m, func() tea.Msg {
			tracks, err := m.db.SearchTracks(msg.query)
			if err != nil {
				return nil
			}
			return TracksLoadedMsg{tracks: tracks}
		}

	case tea.KeyMsg:
		// Global shortcuts
		switch msg.String() {
		case "q", "ctrl+c":
			if m.audioEngine != nil {
				m.audioEngine.Close()
			}
			if m.db != nil {
				m.db.Close()
			}
			return m, tea.Quit
		}

		// Handle library navigation
		switch msg.String() {
		case "tab":
			// Switch between All Songs and Playlists
			if m.librarySection == "songs" {
				m.librarySection = "playlists"
			} else {
				m.librarySection = "songs"
			}
			return m, nil
		case "j", "down":
			if len(m.tracks) > 0 && m.selectedTrackIdx < len(m.tracks)-1 {
				m.selectedTrackIdx++
				m.adjustScroll()
			}
			return m, nil
		case "k", "up":
			if m.selectedTrackIdx > 0 {
				m.selectedTrackIdx--
				m.adjustScroll()
			}
			return m, nil
		case "enter":
			// Load selected track to active deck
			if m.selectedTrackIdx >= 0 && m.selectedTrackIdx < len(m.tracks) {
				track := m.tracks[m.selectedTrackIdx]
				err := m.audioEngine.LoadToDeck(track.Path, m.deckMode)
				if err != nil {
					m.status = fmt.Sprintf("Error loading %s: %v", track.Title, err)
				} else {
					// Update the appropriate deck track
					if m.deckMode == 'A' {
						m.deckATrack = track
					} else {
						m.deckBTrack = track
					}
					m.status = fmt.Sprintf("Loaded to Deck %c: %s - %s", m.deckMode, track.Artist, track.Title)
				}
			}
			return m, nil
		case "/":
			m.searching = true
			m.searchQuery = ""
			return m, nil
		}

		// Handle search input
		if m.searching {
			switch msg.String() {
			case "enter":
				m.searching = false
				return m, func() tea.Msg {
					if m.searchQuery == "" {
						tracks, err := m.db.GetAllTracks()
						if err != nil {
							return nil
						}
						return TracksLoadedMsg{tracks: tracks}
					}
					tracks, err := m.db.SearchTracks(m.searchQuery)
					if err != nil {
						return nil
					}
					return TracksLoadedMsg{tracks: tracks}
				}
			case "esc":
				m.searching = false
				m.searchQuery = ""
				return m, nil
			case "backspace":
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				}
				return m, nil
			default:
				if len(msg.String()) == 1 {
					m.searchQuery += msg.String()
				}
				return m, nil
			}
		}

		// Main shortcuts
		switch msg.String() {
		case "a", "A":
			m.deckMode = 'A'
			m.status = "Deck A mode"
			return m, nil
		case "b", "B":
			m.deckMode = 'B'
			m.status = "Deck B mode"
			return m, nil
		case "s":
			// Scan music directory
			return m, func() tea.Msg {
				homeDir, _ := os.UserHomeDir()
				musicDir := homeDir + "/repositories/music-samples"
				err := m.scanner.ScanDirectory(musicDir)
				if err != nil {
					return nil
				}
				tracks, err := m.db.GetAllTracks()
				if err != nil {
					return nil
				}
				return TracksLoadedMsg{tracks: tracks}
			}
		case "t":
			// Load MP3 test file to active deck
			err := m.audioEngine.LoadToDeck("/tmp/power.mp3", m.deckMode)
			if err != nil {
				m.status = fmt.Sprintf("Error loading file: %v", err)
			} else {
				track := &database.Track{
					Title:  "Power",
					Artist: "Shanghai Doom",
				}
				if m.deckMode == 'A' {
					m.deckATrack = track
				} else {
					m.deckBTrack = track
				}
				m.status = fmt.Sprintf("Loaded test track to Deck %c", m.deckMode)
			}
			return m, nil
		case "x":
			// Stop active deck
			if m.audioEngine != nil {
				err := m.audioEngine.StopDeck(m.deckMode)
				if err != nil {
					m.status = fmt.Sprintf("Error stopping deck %c: %v", m.deckMode, err)
				} else {
					m.status = fmt.Sprintf("Stopped Deck %c", m.deckMode)
				}
			}
			return m, nil
		case " ", "space", "p":
			// Toggle play/pause on active deck (space or 'p')
			if m.audioEngine != nil {
				isLoaded := m.audioEngine.IsLoaded(m.deckMode)
				isPlaying := m.audioEngine.IsPlaying(m.deckMode)

				if !isLoaded {
					m.status = fmt.Sprintf("No track loaded on Deck %c", m.deckMode)
					return m, nil
				}

				if isPlaying {
					err := m.audioEngine.PauseDeck(m.deckMode)
					if err != nil {
						m.status = fmt.Sprintf("Error pausing deck %c: %v", m.deckMode, err)
					} else {
						m.status = fmt.Sprintf("Paused Deck %c", m.deckMode)
					}
				} else {
					err := m.audioEngine.PlayDeck(m.deckMode)
					if err != nil {
						m.status = fmt.Sprintf("Error playing deck %c: %v", m.deckMode, err)
					} else {
						m.status = fmt.Sprintf("Playing Deck %c", m.deckMode)
					}
				}
			} else {
				m.status = "Audio engine not available"
			}
			return m, nil
		case "+", "=":
			// Increase volume
			if m.audioEngine != nil {
				currentVol := m.audioEngine.GetVolume(m.deckMode)
				newVol := currentVol + 0.1
				if newVol > 1.0 {
					newVol = 1.0
				}
				m.audioEngine.SetVolume(m.deckMode, newVol)
				m.status = fmt.Sprintf("Deck %c volume: %.0f%%", m.deckMode, newVol*100)
			}
			return m, nil
		case "-", "_":
			// Decrease volume
			if m.audioEngine != nil {
				currentVol := m.audioEngine.GetVolume(m.deckMode)
				newVol := currentVol - 0.1
				if newVol < 0.0 {
					newVol = 0.0
				}
				m.audioEngine.SetVolume(m.deckMode, newVol)
				m.status = fmt.Sprintf("Deck %c volume: %.0f%%", m.deckMode, newVol*100)
			}
			return m, nil
		case "c":
			// Cue point - jump to beginning
			if m.audioEngine != nil && m.audioEngine.IsLoaded(m.deckMode) {
				// In a real implementation, would jump to cue point
				m.status = fmt.Sprintf("Cue point set on Deck %c", m.deckMode)
			}
			return m, nil
		case "1", "2", "3", "4", "5", "6", "7", "8":
			// Hot cue points
			if m.audioEngine != nil && m.audioEngine.IsLoaded(m.deckMode) {
				cueNum := msg.String()
				m.status = fmt.Sprintf("Hot cue %s triggered on Deck %c", cueNum, m.deckMode)
			}
			return m, nil
		}
	}
	return m, nil
}

func (m *Model) adjustScroll() {
	libraryHeight := m.height/2 - 6 // Bottom half minus some padding
	visibleLines := libraryHeight - 4 // Account for headers and borders

	if m.selectedTrackIdx < m.scrollOffset {
		m.scrollOffset = m.selectedTrackIdx
	} else if m.selectedTrackIdx >= m.scrollOffset+visibleLines {
		m.scrollOffset = m.selectedTrackIdx - visibleLines + 1
	}
}

func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	return m.renderDJLayout()
}

func (m Model) renderDJLayout() string {
	// Calculate layout dimensions
	topHalfHeight := m.height / 2
	bottomHalfHeight := m.height - topHalfHeight

	// === TOP SECTION: DECKS ===
	decksSection := m.renderDecksSection(m.width, topHalfHeight)

	// === BOTTOM SECTION: LIBRARY ===
	librarySection := m.renderLibrarySection(m.width, bottomHalfHeight)

	// Combine sections
	content := lipgloss.JoinVertical(lipgloss.Top, decksSection, librarySection)
	return content
}

func (m Model) renderDecksSection(width, height int) string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(width-2) // Account for terminal borders

	deckWidth := (width - 8) / 2 // Leave space for gap, ensure fits in terminal
	deckHeight := height - 6 // Leave more space for header and status

	// Deck styles - ensure they fit within terminal bounds
	deckAStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 1).
		Width(deckWidth).
		Height(deckHeight)

	deckBStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 1).
		Width(deckWidth).
		Height(deckHeight)

	// Highlight active deck
	if m.deckMode == 'A' {
		deckAStyle = deckAStyle.BorderForeground(lipgloss.Color("#00FF00")).Bold(true)
		deckBStyle = deckBStyle.BorderForeground(lipgloss.Color("#666666"))
	} else {
		deckAStyle = deckAStyle.BorderForeground(lipgloss.Color("#666666"))
		deckBStyle = deckBStyle.BorderForeground(lipgloss.Color("#00FF00")).Bold(true)
	}

	header := headerStyle.Render("MDMA - Music DJ Management Application")

	// Deck content
	deckAContent := m.renderDeckContent('A')
	deckBContent := m.renderDeckContent('B')

	deckA := deckAStyle.Render(deckAContent)
	deckB := deckBStyle.Render(deckBContent)

	// Crossfader section
	crossfader := m.renderCrossfader()

	decks := lipgloss.JoinHorizontal(lipgloss.Center, deckA, crossfader, deckB)

	// Status line
	status := fmt.Sprintf("Active: Deck %c | %s", m.deckMode, m.status)

	return lipgloss.JoinVertical(lipgloss.Left, header, "", decks, "", status)
}

func (m Model) renderDeckContent(deck rune) string {
	var track *database.Track
	if deck == 'A' {
		track = m.deckATrack
	} else {
		track = m.deckBTrack
	}

	var title string
	if deck == 'A' {
		title = "DECK A"
		if m.deckMode == 'A' {
			title += " ●"
		}
	} else {
		title = "DECK B"
		if m.deckMode == 'B' {
			title += " ●"
		}
	}

	if track == nil {
		return fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s\n\n%s",
			title,
			"No track loaded",
			m.renderWaveform(deck, false),
			"[⏸] [Cue] [Sync]",
			m.renderVolumeControls(deck))
	}

	// Track info with metadata
	trackInfo := fmt.Sprintf("%s - %s", track.Artist, track.Title)
	if len(trackInfo) > 30 {
		trackInfo = trackInfo[:27] + "..."
	}

	// Additional metadata
	metadata := ""
	if track.BPM > 0 {
		metadata += fmt.Sprintf("%.0f BPM", track.BPM)
	}
	if track.Key != "" {
		if metadata != "" {
			metadata += " • "
		}
		metadata += track.Key
	}
	if track.Genre != "" {
		if metadata != "" {
			metadata += " • "
		}
		metadata += track.Genre
	}

	// Track progress and time
	timeInfo := m.renderTrackTime(deck)

	// Waveform
	waveform := m.renderWaveform(deck, true)

	// Controls
	controls := m.renderDeckControls(deck)

	// Volume controls
	volume := m.renderVolumeControls(deck)

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s",
		title, trackInfo, metadata, timeInfo, waveform, controls, volume)
}

func (m Model) renderLibrarySection(width, height int) string {
	leftPaneWidth := width / 3
	rightPaneWidth := width - leftPaneWidth - 4 // Account for borders and spacing

	leftPane := m.renderLeftPane(leftPaneWidth, height-2)
	rightPane := m.renderTrackTable(rightPaneWidth, height-2)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, "  ", rightPane)
}

func (m Model) renderLeftPane(width, height int) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#04B575")).
		Padding(1, 1).
		Width(width).
		Height(height)

	// All Songs section
	songsStyle := lipgloss.NewStyle().Bold(m.librarySection == "songs")
	songsText := songsStyle.Render("All Songs")
	if m.librarySection == "songs" {
		songsText = lipgloss.NewStyle().
			Background(lipgloss.Color("#7D56F4")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).
			Render("All Songs")
	}

	// Playlists section
	playlistsStyle := lipgloss.NewStyle().Bold(m.librarySection == "playlists")
	playlistsText := playlistsStyle.Render("Playlists")
	if m.librarySection == "playlists" {
		playlistsText = lipgloss.NewStyle().
			Background(lipgloss.Color("#7D56F4")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).
			Render("Playlists")
	}

	content := fmt.Sprintf("%s\n\n%s\n  My Playlist\n  DJ Set 1\n  Favorites", songsText, playlistsText)

	return style.Render(content)
}

func (m Model) renderTrackTable(width, height int) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#04B575")).
		Padding(1, 1).
		Width(width).
		Height(height)

	var header string
	if m.searching {
		header = fmt.Sprintf("Search: %s_", m.searchQuery)
	} else {
		header = fmt.Sprintf("%s (%d tracks)", m.librarySection, len(m.tracks))
	}

	if len(m.tracks) == 0 {
		content := fmt.Sprintf("%s\n\nNo tracks found.\nPress 's' to scan music directory.", header)
		return style.Render(content)
	}

	// Table header
	tableHeader := fmt.Sprintf("%-25s %-20s %-6s %-15s", "Name", "Artist", "BPM", "Comment")
	separator := strings.Repeat("─", width-4)

	lines := []string{header, "", tableHeader, separator}

	// Calculate visible area
	visibleLines := height - 8 // Account for header, border, padding
	start := m.scrollOffset
	end := start + visibleLines
	if end > len(m.tracks) {
		end = len(m.tracks)
	}

	// Track rows
	for i := start; i < end; i++ {
		track := m.tracks[i]
		line := m.formatTrackTableRow(track, width-6)

		if i == m.selectedTrackIdx {
			line = lipgloss.NewStyle().
				Background(lipgloss.Color("#7D56F4")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Width(width-6).
				Render(line)
		}

		lines = append(lines, line)
	}

	// Add scroll indicator
	if len(m.tracks) > visibleLines {
		scrollInfo := fmt.Sprintf("(%d-%d of %d)", start+1, end, len(m.tracks))
		lines = append(lines, "", scrollInfo)
	}

	// Help text
	help := "[j/k] Navigate [Enter] Load [Tab] Section [/] Search [A/B] Deck [+/-] Volume [c] Cue [1-8] Hot Cues [q] Quit"
	lines = append(lines, "", help)

	content := strings.Join(lines, "\n")
	return style.Render(content)
}

func (m Model) formatTrackTableRow(track *database.Track, width int) string {
	name := track.Title
	if len(name) > 24 {
		name = name[:21] + "..."
	}

	artist := track.Artist
	if len(artist) > 19 {
		artist = artist[:16] + "..."
	}

	bpm := ""
	if track.BPM > 0 {
		bpm = fmt.Sprintf("%.0f", track.BPM)
	}

	comment := track.Notes
	if len(comment) > 14 {
		comment = comment[:11] + "..."
	}

	return fmt.Sprintf("%-25s %-20s %-6s %-15s", name, artist, bpm, comment)
}

func (m Model) renderWaveform(deck rune, hasTrack bool) string {
	if !hasTrack {
		return "░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░"
	}

	// Simulate a waveform with different heights
	waveform := ""
	wavePattern := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

	// Create a simple waveform pattern (in real implementation would use actual audio data)
	for i := 0; i < 30; i++ {
		height := (i*3 + int(deck)*7) % len(wavePattern)
		if deck == 'A' {
			waveform += lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF88")).Render(wavePattern[height])
		} else {
			waveform += lipgloss.NewStyle().Foreground(lipgloss.Color("#0088FF")).Render(wavePattern[height])
		}
	}
	return waveform
}

func (m Model) renderTrackTime(deck rune) string {
	// In a real implementation, would track actual playback position
	// For now, show placeholder times
	return "1:23 / 4:56  ████████░░░░░░░░"
}

func (m Model) renderDeckControls(deck rune) string {
	isPlaying := m.audioEngine != nil && m.audioEngine.IsPlaying(deck)

	var playButton string
	if isPlaying {
		playButton = lipgloss.NewStyle().
			Background(lipgloss.Color("#00FF00")).
			Foreground(lipgloss.Color("#000000")).
			Render("[▶]")
	} else {
		playButton = "[⏸]"
	}

	cueButton := lipgloss.NewStyle().
		Background(lipgloss.Color("#FF8800")).
		Foreground(lipgloss.Color("#000000")).
		Render("[Cue]")

	syncButton := "[Sync]"

	// Hot cue points
	hotCues := ""
	for i := 1; i <= 8; i++ {
		if i <= 4 {
			hotCues += fmt.Sprintf("[%d]", i)
		} else if i == 5 {
			hotCues += "\n"
		}
		if i > 4 {
			hotCues += fmt.Sprintf("[%d]", i)
		}
	}

	return fmt.Sprintf("%s %s %s\n\nHot Cues:\n%s", playButton, cueButton, syncButton, hotCues)
}

func (m Model) renderVolumeControls(deck rune) string {
	volume := 1.0
	if m.audioEngine != nil {
		volume = m.audioEngine.GetVolume(deck)
	}

	// Create visual volume slider
	totalBars := 10
	filledBars := int(volume * float64(totalBars))

	volumeBar := ""
	for i := 0; i < totalBars; i++ {
		if i < filledBars {
			volumeBar += "█"
		} else {
			volumeBar += "░"
		}
	}

	bpm := "---"
	var track *database.Track
	if deck == 'A' {
		track = m.deckATrack
	} else {
		track = m.deckBTrack
	}
	if track != nil && track.BPM > 0 {
		bpm = fmt.Sprintf("%.0f", track.BPM)
	}

	return fmt.Sprintf("Vol: %s  BPM: %s  Gain: ±0", volumeBar, bpm)
}

func (m Model) renderCrossfader() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF6B9D")).
		Padding(1, 1).
		Width(20).
		Align(lipgloss.Center)

	// Crossfader position (center for now)
	crossfader := "A ◄─────●─────► B"

	// BPM sync display
	bpmA := "---"
	bpmB := "---"
	if m.deckATrack != nil && m.deckATrack.BPM > 0 {
		bpmA = fmt.Sprintf("%.0f", m.deckATrack.BPM)
	}
	if m.deckBTrack != nil && m.deckBTrack.BPM > 0 {
		bpmB = fmt.Sprintf("%.0f", m.deckBTrack.BPM)
	}

	// Sync status
	syncStatus := ""
	if bpmA != "---" && bpmB != "---" {
		diffBPM := ""
		if m.deckATrack != nil && m.deckBTrack != nil {
			diff := m.deckATrack.BPM - m.deckBTrack.BPM
			if diff > 0 {
				diffBPM = fmt.Sprintf("+%.1f", diff)
			} else if diff < 0 {
				diffBPM = fmt.Sprintf("%.1f", diff)
			} else {
				diffBPM = "SYNC"
				syncStatus = lipgloss.NewStyle().
					Background(lipgloss.Color("#00FF00")).
					Foreground(lipgloss.Color("#000000")).
					Render("SYNC")
			}
		}
		if syncStatus == "" {
			syncStatus = diffBPM
		}
	}

	content := fmt.Sprintf("CROSSFADER\n\n%s\n\n%s ↔ %s\n%s", crossfader, bpmA, bpmB, syncStatus)

	return style.Render(content)
}