# Detailed Design Document

## System Overview
The solution comprises three major components: (1) a self-hosted media server exposing the MP3 library and metadata APIs, (2) client interfaces (mobile, desktop, TUI) consuming those APIs, and (3) a shared metadata service ensuring consistency for playlists, tags, and cue points.

## Media Server Layer
- Deploy Navidrome or Jellyfin as the baseline server to catalogue audio files, stream tracks, and expose REST/GraphQL endpoints for library operations.
- Extend via plugins/webhooks to handle custom tag schema and cue point storage in a companion metadata service (e.g., Postgres + Hasura).
- Secure access through WireGuard/Tailscale; store secrets in environment-managed vault.

## Metadata Service
- Postgres schema: `tracks`, `playlists`, `tags`, `track_tags`, `cue_points`, `sessions`.
- Provide GraphQL/REST API for CRUD operations; include optimistic locking for conflict resolution.
- Background jobs sync metadata back to server-side files (ID3) to keep local copies authoritative.

## Client Applications
- **Mobile**: Leverage existing open-source clients where possible (Jellyfin mobile, Symfonium) configured to point at the server; supplement with a thin companion app/webview for tag/cue editing.
- **Desktop Web**: Use server-provided UI for browsing; integrate custom metadata management panels.
- **TUI DJ App** (primary focus):
  - Built in Go using `tcell` or `bubbletea` for UI, `oto` or `portaudio` bindings for playback.
  - Architecture: `pkg/ui`, `pkg/audio`, `pkg/library`, `pkg/session`; command bus for keyboard shortcuts.
  - Waveform pre-processing service generates waveform data and stores alongside tracks.
  - Provide dual-deck layout with crossfader, tempo adjust, cue management, track browser.

## Data Flow
1. User selects track in TUI; library module fetches metadata and waveform from API.
2. Audio module streams track via server’s transcoding endpoint, buffering locally.
3. UI module renders waveform and handles keyboard events to cue, play, and adjust tempo.
4. Metadata edits pushed to metadata service; background sync updates ID3 tags.

## Non-Functional Requirements
- Latency under 50ms for audio adjustments; requires local caching and preloading.
- Offline cache for up to 100 tracks with automatic invalidation.
- Unit and integration tests for metadata operations; signal handling tests for playback reliability.
