# Implementation Plan

## Phase 1: Foundations
1. Stand up media server (Navidrome/Jellyfin) in local docker-compose; validate remote access via WireGuard.
2. Design Postgres schema for metadata service; scaffold Go backend with ent or sqlc.
3. Implement REST/GraphQL endpoints for tracks, playlists, tags, cue points with optimistic concurrency.
4. Create waveform extraction worker using ffmpeg + Go job runner; store data in object store (S3-compatible).

## Phase 2: Metadata Sync & Clients
1. Develop background synchronization that maps metadata updates to ID3 tags and pushes to server library.
2. Integrate mobile/desktop clients with metadata APIs; add custom UI for tags and cue points (web panel or PWA).
3. Establish authentication flow (OpenID Connect) shared across all clients.

## Phase 3: TUI Application
1. Bootstrap Go module (`cmd/mdma`) with `bubbletea` UI skeleton; set up `pkg` directories for audio, UI, library.
2. Implement library browser, search, filtering, and tagging flows using API client.
3. Build dual-deck playback: streaming buffers, tempo adjustment, crossfader, waveform rendering.
4. Add cue point management (set, jump, label), keyboard shortcuts, and session persistence.

## Phase 4: Quality & Deployment
1. Write unit tests for metadata service and audio abstractions; create integration tests using Docker Compose.
2. Set up CI pipeline (GitHub Actions) for linting (`golangci-lint`), formatting (`gofmt`), tests, and build artifacts.
3. Provide deployment scripts (Terraform/Ansible) for server + metadata stack; document monitoring (Prometheus/Grafana).
