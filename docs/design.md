# MDMA - Music DJ Management Application Design Document

## Overview

MDMA is a terminal-based music management and DJ preparation tool designed for track auditioning, set preparation, and library management. It focuses on workflow efficiency for amateur DJs who need to organize music, create transitions, and prepare sets.

## System Architecture

### Components

1. **Music Server** - Navidrome (lightweight, excellent API, Subsonic compatible)
2. **TUI Client** - Go application using Bubble Tea framework
3. **Audio Engine** - PortAudio + Go bindings for cross-platform audio
4. **Storage** - Local SQLite for metadata, tags, and cue points

### Technology Stack

- **Language**: Go (performance, deployment simplicity, strong TUI ecosystem)
- **TUI Framework**: Bubble Tea + Lipgloss (modern, actively developed)
- **Audio**: PortAudio (mature, cross-platform, low-latency capable)
- **HTTP Client**: Standard library + retries for server communication
- **Database**: SQLite for local metadata and offline capability
- **Config**: YAML configuration files

## Core Features

### 1. Library Management

**Music Server Integration**
- Connect to Navidrome server via Subsonic API
- Browse/search entire music library
- Sync metadata locally for offline browsing
- Support for multiple servers (home, cloud backup)

**Local Metadata Enhancement**
- Custom tagging system (Energy, Mood, Key, etc.)
- Star ratings and personal notes
- Track grouping and relationship mapping
- Import/export tag data

**Playlist Management**
- Create/edit playlists locally and sync to server
- Smart playlists based on tags/metadata
- Set preparation playlists with transition notes

### 2. Dual-Deck Player

**Basic Playback**
- Two independent audio channels
- Play/pause/stop controls per deck
- Volume faders (0-100%)
- Crossfader for mixing between decks
- Tempo adjustment (+/- 20% with pitch preservation)

**Track Loading**
- Drag tracks from library to decks
- Queue system for next tracks
- Auto-load next track option
- Track history and recently played

**Visual Feedback**
- ASCII waveform display (simplified, pre-computed)
- Current position indicator
- BPM display (if available in metadata)
- Track information (title, artist, key, energy)

### 3. Cue Point System

**Cue Management**
- 8 cue points per track (hot cues)
- Text labels for each cue point
- Visual markers on waveform
- Jump to cue functionality

**Common Cue Types**
- Intro/outro markers
- Drop points
- Breakdown sections
- Mix-in/mix-out points

### 4. Transition Planning

**Track Analysis**
- BPM detection and display
- Key detection (if available)
- Energy level assessment
- Compatible track suggestions

**Mix Notes**
- Per-track mixing notes
- Transition timing suggestions
- Equipment-specific notes (loops, effects)

## User Interface Design

### Layout Structure

```
┌─ MDMA v1.0 ─────────────────────────────────────────────────────────────┐
│                                                                         │
│ ┌─ DECK A ─────────────────┐  ┌─ DECK B ─────────────────┐             │
│ │ Artist - Track Name      │  │ Artist - Track Name      │             │
│ │ ▓▓▓▓▓▓▓▓░░░░░░░░░░░░ 60% │  │ ░░░░░░░░▓▓▓▓▓▓▓▓▓▓▓▓ 80% │             │
│ │ [▶] [⏸] [⏹] BPM: 128   │  │ [▶] [⏸] [⏹] BPM: 126   │             │
│ │ Tempo: +2% Vol: ███░░   │  │ Tempo: -1% Vol: ████░   │             │
│ └─────────────────────────┘  └─────────────────────────┘             │
│                                                                         │
│ ┌─ CROSSFADER ──────────────────────────────────────────────────────┐   │
│ │                     A ◄──●──► B                                   │   │
│ └───────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│ ┌─ LIBRARY ──────────────────┐ ┌─ PLAYLIST ─────────────────────────┐   │
│ │ > Electronic               │ │ > Saturday Night Set               │   │
│ │   > House                  │ │   Track 1 - Artist A               │   │
│ │     Track Name - Artist    │ │   Track 2 - Artist B               │   │
│ │     Another Track          │ │   Track 3 - Artist C               │   │
│ │   > Techno                 │ │                                    │   │
│ │ > Hip Hop                  │ │ [a] Add Track [r] Remove [e] Edit  │   │
│ └────────────────────────────┘ └────────────────────────────────────┘   │
│                                                                         │
│ [l] Library [p] Playlists [s] Search [t] Tags [q] Quit                  │
└─────────────────────────────────────────────────────────────────────────┘
```

### Key Bindings (lazygit-inspired)

**Global Navigation**
- `l` - Library view
- `p` - Playlists view
- `s` - Search
- `t` - Tag editor
- `q` - Quit

**Library Management**
- `j/k` - Navigate up/down
- `enter` - Load to active deck
- `1/2` - Load to specific deck
- `space` - Preview track
- `a` - Add to playlist
- `tab` - Switch between library/playlist panes

**Player Controls**
- `z/x` - Play/pause deck A/B
- `c/v` - Stop deck A/B
- `f/g` - Tempo adjust deck A/B
- `r/t` - Volume adjust deck A/B
- `b` - Toggle crossfader mode

**Cue Points**
- `F1-F8` - Set/jump to cue points
- `shift+F1-F8` - Clear cue points

## Technical Implementation Details

### Audio Pipeline

```
[File] → [Decoder] → [Resampler] → [Effects] → [Mixer] → [Output]
                                      ↑
                                [Tempo/Pitch]
```

**File Support**
- MP3, FLAC, AAC, OGG (via Go audio libraries)
- Streaming from HTTP endpoints (Navidrome)
- Local file access with metadata reading

**Audio Processing**
- Real-time tempo adjustment without pitch change
- Volume control per channel
- Simple crossfading algorithm
- Basic EQ (optional future feature)

### Data Management

**Local Database Schema**
```sql
-- Custom metadata storage
tracks (
    id INTEGER PRIMARY KEY,
    server_id TEXT,
    local_path TEXT,
    custom_tags TEXT, -- JSON
    energy_level INTEGER,
    cue_points TEXT, -- JSON
    notes TEXT,
    rating INTEGER
);

playlists (
    id INTEGER PRIMARY KEY,
    name TEXT,
    tracks TEXT, -- JSON array
    created_at TIMESTAMP
);
```

**Configuration Management**
```yaml
# ~/.config/mdma/config.yaml
servers:
  - name: "home"
    url: "http://192.168.1.100:4533"
    username: "user"
    password: "pass"

audio:
  buffer_size: 512
  sample_rate: 44100
  output_device: "default"

ui:
  theme: "dark"
  waveform_detail: "medium"
```

### Server Communication

**Navidrome API Integration**
- Authentication via Subsonic protocol
- Streaming URL generation
- Metadata synchronization
- Playlist sync (bidirectional)

**Offline Capabilities**
- Cache frequently accessed metadata
- Offline playlist editing
- Sync changes when connection restored

## Performance Considerations

**Memory Management**
- Stream audio data rather than loading entire files
- LRU cache for frequently accessed tracks
- Lazy loading of library metadata

**CPU Usage**
- Efficient waveform generation (pre-computed when possible)
- Minimal UI redraws
- Background metadata processing

**Network Usage**
- Intelligent prefetching of likely-to-be-played tracks
- Compressed metadata sync
- Resume interrupted downloads

## Future Extensions

**Phase 2 Features**
- BPM detection and sync
- Key detection and harmonic mixing
- Effect processing (filters, reverb)
- Recording/streaming output

**Integration Possibilities**
- Spotify playlist import (via Web API)
- Last.fm scrobbling
- Hardware controller support (MIDI)
- Plugin system for custom effects

This design balances the simplicity needed for a TUI with the functionality required for effective DJ preparation and music management.