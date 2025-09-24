# MDMA Implementation Plan

## Development Phases

### Phase 1: Core Foundation (Weeks 1-3)

**1.1 Project Setup**
- Initialize Go module with proper structure
- Set up development environment and tooling
- Configure CI/CD pipeline (GitHub Actions)
- Create basic project documentation

**1.2 Audio Foundation**
- Research and integrate PortAudio Go bindings
- Implement basic audio playback for single files
- Create audio pipeline abstraction
- Test cross-platform audio support (macOS primary)

**1.3 TUI Framework**
- Set up Bubble Tea application structure
- Implement basic layout and navigation
- Create component system for reusable UI elements
- Implement keyboard event handling

**1.4 Configuration System**
- Implement YAML configuration loading
- Create configuration validation
- Set up default configuration generation
- Add configuration hot-reloading

### Phase 2: Library Management (Weeks 4-6)

**2.1 Local Database**
- Set up SQLite integration
- Implement database schema and migrations
- Create data access layer (DAL)
- Add database connection pooling

**2.2 File System Integration**
- Implement local music file scanning
- Add metadata extraction (ID3, FLAC tags)
- Create file watcher for library changes
- Implement library indexing and search

**2.3 Basic Library UI**
- Create tree view for music library browsing
- Implement search and filtering
- Add track selection and preview
- Create basic playlist management UI

### Phase 3: Navidrome Integration (Weeks 7-8)

**3.1 Subsonic API Client**
- Implement Subsonic API authentication
- Create API client for basic operations
- Add error handling and retry logic
- Implement streaming URL generation

**3.2 Server Integration**
- Add server configuration management
- Implement metadata synchronization
- Create offline/online mode switching
- Add playlist synchronization

### Phase 4: Dual-Deck Player (Weeks 9-11)

**4.1 Audio Engine Enhancement**
- Implement dual-channel audio mixing
- Add volume and tempo controls
- Create crossfader functionality
- Implement audio synchronization

**4.2 Player UI Components**
- Create deck display components
- Implement waveform visualization (ASCII)
- Add transport controls (play/pause/stop)
- Create crossfader interface

**4.3 Track Loading System**
- Implement drag-and-drop to decks
- Add queue management
- Create auto-advance functionality
- Implement track history

### Phase 5: Advanced Features (Weeks 12-14)

**5.1 Cue Point System**
- Implement cue point storage and retrieval
- Add cue point UI markers
- Create hot cue keyboard shortcuts
- Implement cue point management

**5.2 Custom Tagging**
- Create tag management system
- Implement tag editor UI
- Add batch tagging operations
- Create tag-based search and filtering

**5.3 Playlist Enhancement**
- Add smart playlist creation
- Implement playlist export/import
- Create set preparation features
- Add transition planning tools

### Phase 6: Polish and Testing (Weeks 15-16)

**6.1 Performance Optimization**
- Profile and optimize audio performance
- Implement efficient UI rendering
- Add memory usage optimization
- Create performance benchmarks

**6.2 Testing and Documentation**
- Add comprehensive unit tests
- Create integration tests
- Write user documentation
- Add troubleshooting guides

## Technical Implementation Details

### Project Structure

```
mdma/
├── cmd/
│   └── mdma/
│       └── main.go
├── internal/
│   ├── audio/
│   │   ├── engine.go
│   │   ├── decoder.go
│   │   └── mixer.go
│   ├── config/
│   │   ├── config.go
│   │   └── defaults.go
│   ├── database/
│   │   ├── db.go
│   │   ├── migrations/
│   │   └── models/
│   ├── library/
│   │   ├── scanner.go
│   │   ├── metadata.go
│   │   └── searcher.go
│   ├── server/
│   │   ├── subsonic.go
│   │   └── navidrome.go
│   ├── ui/
│   │   ├── app.go
│   │   ├── components/
│   │   ├── pages/
│   │   └── styles/
│   └── utils/
├── pkg/
│   └── mdma/
├── configs/
├── docs/
├── scripts/
└── tests/
```

### Key Dependencies

```go
// go.mod
module github.com/bradstewart/mdma

go 1.21

require (
    github.com/charmbracelet/bubbletea v0.24.2
    github.com/charmbracelet/lipgloss v0.9.1
    github.com/gordonklaus/portaudio v0.0.0-20230709114228-aafa478834f5
    github.com/hajimehoshi/go-mp3 v0.3.4
    github.com/jfreymuth/oggvorbis v1.0.5
    github.com/mattn/go-sqlite3 v1.14.17
    github.com/spf13/viper v1.16.0
    github.com/dhowden/tag v0.0.0-20230630033851-978a0926ee25
    gopkg.in/yaml.v3 v3.0.1
)
```

### Development Milestones

**Milestone 1 (End of Week 3): Audio Playback**
- Single file playback working
- Basic TUI layout rendered
- Configuration system functional

**Milestone 2 (End of Week 6): Local Library**
- Local music library scanning
- Basic library browsing UI
- Search and filtering working

**Milestone 3 (End of Week 8): Server Integration**
- Navidrome connection established
- Remote library browsing
- Track streaming functional

**Milestone 4 (End of Week 11): Dual-Deck System**
- Two-deck audio mixing
- Basic waveform display
- Transport controls working

**Milestone 5 (End of Week 14): Feature Complete**
- Cue points implemented
- Custom tagging system
- Advanced playlist management

**Milestone 6 (End of Week 16): Release Ready**
- Performance optimized
- Thoroughly tested
- Documentation complete

## Risk Mitigation

### Technical Risks

**Audio Latency Issues**
- *Risk*: Audio dropouts or high latency on macOS
- *Mitigation*: Use Core Audio directly if PortAudio insufficient
- *Fallback*: Reduce buffer size, implement audio prioritization

**Memory Usage**
- *Risk*: High memory usage with large libraries
- *Mitigation*: Implement streaming and caching strategies
- *Fallback*: Add memory usage monitoring and limits

**Cross-Platform Compatibility**
- *Risk*: Audio issues on different platforms
- *Mitigation*: Test early and often on target platforms
- *Fallback*: Platform-specific audio backends

### Development Risks

**Scope Creep**
- *Risk*: Adding too many features early
- *Mitigation*: Strict phase adherence, MVP focus
- *Fallback*: Cut non-essential features

**Third-Party Dependencies**
- *Risk*: Breaking changes in Bubble Tea or audio libraries
- *Mitigation*: Pin versions, maintain compatibility layer
- *Fallback*: Fork dependencies if necessary

## Testing Strategy

### Unit Testing
- All audio processing functions
- Configuration loading and validation
- Database operations
- UI component logic

### Integration Testing
- Audio pipeline end-to-end
- Navidrome API integration
- File system operations
- UI interaction flows

### Performance Testing
- Audio latency measurements
- Memory usage profiling
- UI responsiveness under load
- Large library handling

### Manual Testing
- Real DJ workflow scenarios
- Cross-platform compatibility
- Error condition handling
- User experience validation

## Deployment Strategy

### Development Environment
- Go 1.21+ required
- PortAudio system dependencies
- SQLite for development database
- Make-based build system

### Distribution
- Single binary distribution
- Platform-specific builds (macOS, Linux)
- Homebrew formula for macOS
- GitHub Releases for distribution

### Documentation
- README with quick start guide
- Installation instructions per platform
- Configuration reference
- Troubleshooting guide
- API documentation for extensibility

This implementation plan prioritizes getting a working audio foundation early, then building up the UI and features incrementally. Each phase delivers working functionality that can be tested and validated before moving to the next phase.