# High-Level Observations

## Glaring Holes & Pitfalls
- **Server selection ambiguity**: unclear whether Jellyfin or another server offers required two-way metadata sync, tag management, and multi-device writes; off-the-shelf options may resist custom tagging workflows.
- **Remote access & storage**: no plan for securely exposing the home MP3 library to outside devices; bandwidth, latency, and NAT traversal issues remain unsolved.
- **Mobile integration**: expectation of Spotify + local library within one app is risky given Spotify API rate limits and authentication constraints; native apps seldom allow third-party local playback.
- **Audio playback complexity**: real-time mixing, tempo sync, waveform rendering, and cue points require low-latency DSP; go/python ecosystem for cross-platform TUI audio is thin.
- **Metadata consistency**: simultaneous edits from multiple devices demand conflict resolution strategy and shared schema for tags/playlists.

## Questions
- What is the authoritative source for metadata — the home server or a cloud sync service?
- Is lossy streaming acceptable for remote playback, or is lossless required?
- How critical is offline access on mobile devices?
- Are there existing DJ controllers or MIDI devices that should integrate with the TUI?

## Proposed Next Steps
1. Evaluate self-hosted media servers (Jellyfin, Navidrome, Funkwhale) for metadata write support and plugin capacity.
2. Prototype secure remote access (WireGuard or Tailscale) to expose the media server.
3. Spike on waveform extraction + playback control libraries in Go (e.g., Ebiten, PortAudio bindings) and Python (e.g., pyaudio, librosa + textual UI).
4. Define metadata schema and conflict-resolution rules before UI work.
