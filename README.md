# My DJ Mixing App
(lol, naming is... hard?). A terminal-based audio player and library manager for "DJs".

## Goals

This is primarily an experiment and learning exercise.
- Golang. Never used it, want to use it.
- TUIs. Never built one! I also don't spend enough time in the terminal. Maybe this finally makes me learn neovim.

It might turn into something functional that scratches my own itch.

## Functionality

I did mention scratching my own itch... I am a (very) amateur DJ. I mostly play for myself, sometimes a few friends.

I listen to music all daylong whilst writing code for money. I often hear two songs in a row that blend well and think "yo I want to remember those two go together" without breaking the rest of my flow.

I also find myself wanting:
- a better tagging system then what Apple Music/iTunes offers (yes I just use Apple Music.app to manage my mp3s, it was just.... there).
- the ability to construct playlists (and sub-playlists) when preparing mixes -- again in a low-friction, focus-retaining way.
    - arrange songs in an order
    - mark two as "good to together"
    - have a separate list of songs that might work, but aren't yet arranged
    - extra cool: song recommender
- the ability to play two songs at once with some super basic mixing controls (bpm, 3-band EQ, volume, possibly some simple loop controls)

You can do most (all?) of these things in Rekordbox (or Traktor or Mixxx or whatever your DJ software of choice is), _however_ launching those applications is massively focus-destroying (and slow). If I have Rekordbox running, I get no work done. Also the learning/experiment part of this.

### Non-Features (at least not yet)

- File management (e.g. importing mp3 files and organizing them). May add in the future, but lot's of existing tools to do this, so we're assuming your files are reasonably organized on disk already (personally I just drop stuff in Apple Music and let it copy things around).

### Future Ideas

- Streaming service integration. Would be great to tag songs, mark two as "good together", etc from Spotify and other services before I actually go get the song from beatport.
- Song recommender. Dig up stuff that would sound good together.
- Relatedly, song clusterer (for lack of a better word). Basically auto-categorization based on actual audio features. I listen to a lot of weird bass music, and asking the tool for songs with laser noises would be amazing.

### Inspiration

I thought what I wanted didn't really exist in any form for a while, until I found these two things:

[Beatport DJ](https://dj.beatport.com/) looks _very_ much like what I imagine for the player/mixer side of this.
[Lexicon](https://www.lexicondj.com/) looks _very_ much like what I imagine for the library management/tagging/grouping side of this.

But: I don't have a Beatport Streaming subscription, I don't want to pay a monthly subscription for Lexicon (to be fair, it looks like great software, and great software is worth paying for, I just need to keep this a "cheap" hobby), terminal-based tools are less focus disturbing, oh and I want to learn some stuff. It's been too long since I hacked on software just cause I can.

## Status

🎉 **Working prototype!** Full-featured DJ software with professional interface:

**🎛️ Dual-Deck Interface:**
- **Waveform displays** with visual audio representation
- **Track metadata** showing Artist, Title, BPM, Key, Genre
- **Visual volume controls** with real-time level indicators
- **BPM sync detection** and crossfader with tempo comparison
- **Hot cue points** (1-8) for instant track navigation
- **Play/pause/cue controls** with visual feedback

**📚 Music Library:**
- **Professional track table** (Name, Artist, BPM, Comments)
- **All Songs and Playlists** sections with tab navigation
- **Real-time search** and filtering
- **Database-backed metadata** with automatic scanning

**🎵 Audio Support:**
- **MP3 and M4A/AAC** playback via FFmpeg integration
- **Dual-channel audio engine** with independent deck control
- **Volume control** with visual feedback

## Requirements

- **Go 1.24+**
- **FFmpeg** (for M4A/AAC support)
- **Audio hardware** (for actual playback)

## Quick Start

```bash
# Build and run
go build -o mdma cmd/mdma/main.go
./mdma

# Controls
## Deck Controls
[A/B]     Switch active deck
[Space/p] Play/pause active deck
[x]       Stop active deck
[c]       Cue point (return to start)
[1-8]     Hot cue points
[+/-]     Volume up/down
[t]       Load test track

## Library Navigation
[j/k]     Navigate track list
[Enter]   Load selected track to active deck
[Tab]     Switch between All Songs/Playlists
[/]       Search tracks
[s]       Scan music directory
[q]       Quit
```

## Plan

- [ ] Sketch out the UI.
- [ ] Do some research into Go, in general. Project setup, machine configs, etc, etc.
- [ ] Research Go TUI libraries.
- [ ] Research Go audio libraries.
- [ ] Plan to plan?
