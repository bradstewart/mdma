# Overall System

This is a high-level overview of the system we want to build (as of 2025-09-23).

## Background

I am a (very) amateur DJ. I mostly play for myself, sometimes a few friends.

I listen to music all daylong whilst writing code for money. I often hear two songs in arow that blend well and think "yo I want to remember those two go together" without
breaking the rest of my flow.

I primarily use Spotify for listening, and iTunes/Apple Music (via Music.app, but none of
its streaming features) for my mp3 collection.

I'm also (now) in an office, away from my mp3 collection and often on a different laptop.

## What do we want this to do?

- Access my mp3 library (stored at some device in my house) from my Android phone, office
  computer, or any other device.
- Be able to add songs, manage playlists, etc across all devices.
- Be able to manage tags and other metadata across all devices.

## Proposed Design (Super High Level)

- Some sort of music server, ideally an off-the-shelf one like Jellyfin? I don't know what
  my options are here, really.
- A mobile app that can interact with said server. It'd be super nice if I could access
  Spotify in here too, but that's a nice to have. Also preferably pre-existing apps.
- A TUI for managing and playing music -- this is the main thing we want to build here,
  and will be the focu of the rest of this. document

### TUI

Think super-simplified DJ software, running in a terminal. Keyboard-driven shortcuts for
most common actions. Use lazygit as an example here.

#### Library Management

I want standard music library management, playlists, sort/filter, browsing, all
that good stuff. It should also support:

- Custom tagging. I want to have a defined list of tags that I can quickly apply

#### Player

Beatport's web player is a good example of what I want here.

A "play two songs at once" mode:

- two super simple "decks" at the top of the screen, play/pause buttons, tempo controls,
  volume faders.
- Waveform display
- A simple cue point system to mark certain points in the song and attach some text

#### Implementation

Use golang or python for the TUI -- i want this to be stable and fast, but also easyto
read and hack on. Consider what underlying libraries are available to make this work, both
for the TUI itself and for the music mixing and playback.

## Prompt

Plan first. Think hard through the above. Point out any glaring holes in the design, any
likely pitfalls, etc. Ask any questions you might have. Then suggest a plan for how to move forward.

Then, write a detailed design document.

Then, write a detailed implementation plan.
