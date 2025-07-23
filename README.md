[![Tests](https://github.com/chriserin/seq/actions/workflows/test.yml/badge.svg)](https://github.com/chriserin/seq/actions/workflows/test.yml)

# SEQ

A powerful MIDI sequencer designed for the command line

## Overview

SEQ is a terminal-based MIDI sequencer that brings midi sequencer capabilities to your CLI. Built with Go and designed for efficiency, SEQ offers rapid beat creation, complex arrangements, and advanced pattern manipulation through an intuitive keyboard-driven interface.

## Key Features

- **Rapid Beat Creation**: Create drum patterns in just a few keystrokes using pattern mode
- **Advanced Overlays**: Add mathematical variations to sequences with the overlay system
- **Flexible Arrangements**: Structure songs with sections, parts, and groups
- **Real-time Manipulation**: Modify patterns, accents, gates, and timing while playing
- **MIDI Integration**: Full MIDI output support for hardware and software instruments
- **Vim-inspired Navigation**: Familiar key bindings for efficient workflow

## Install from package

Pre-built packages for macOS and Linux are found on the [Releases](https://github.com/chriserin/seq/releases) page.

## Install from Source

### Dependencies MACOS

```
# Lua 5.4
brew install lua@5.4-dev
```

### Dependencies Linux

```
sudo apt-get install liblua5.4-dev
sudo apt-get install libasound2-dev
```

```bash
# Clone the repository
git clone https://github.com/chriserin/seq.git
cd seq

# Build with Lua support
go build -tags lua54 -o seq

# Run
./seq
```

### Requirements

- Go 1.24+
- MIDI output device (hardware or software)

## Quick Start

1. **Launch SEQ**: `./seq`
2. **Create a beat**: Move cursor with `hjkl`, press `1` for notes on every beat
3. **Play**: Press `Space` to play/stop
4. **Save**: `Ctrl+s` to save your sequence

### Basic Beat Creation Example

```
BDK┤▧       ▧       ▧       ▧         
SN │    ▧       ▧       ▧       ▧
H1 │▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧▧
```

Created with just 6 keystrokes:
- Bass drum: cursor on BD line, press `8` (note every 8 beats)
- Snare: `j` (down), `4` then `8` (note every 8 beats starting at beat 5)
- Hi-hat: `j` (down), `1` (note every beat)

## Core Concepts

### Pattern Creation

**Pattern Mode** enables rapid beat creation:
- **Numbers (1-9)**: Add/remove notes every N beats
- **Pattern Mode - Value**: `na` (accent), `nw` (wait), `ng` (gate), `nr` (ratchet)
- **Shift+Numbers**: Increase values, Numbers alone decrease values

> 📖 **Detailed Guide**: [Beat Creation Documentation](docs/beat-creation.md)

### Note Attributes

Each note has four modifiable properties:

- **Accent** (`A`/`a`): MIDI velocity (8 levels)
- **Gate** (`G`/`g`): Note duration (8 levels + long gates with `E`/`e`)
- **Ratchet** (`R`/`r`): Multiple hits per beat (8 levels)
- **Wait** (`W`/`w`): Swing timing delay (8 levels)

> 📖 **Detailed Guide**: [Note Alteration Documentation](docs/note-alteration.md)

### Key Cycles and Arrangements

**Key Cycles** track sequence iterations and drive both arrangements and overlays:
- Key line (marked with `K`) determines when cycles advance
- Arrangements use key cycles to switch between sections
- Overlays use key cycles to determine when variations apply

**Arrangements** structure your songs:
- **Sections**: Containers that reference parts with playback attributes
- **Parts**: The actual musical content
- **Groups**: Collections of sections that can repeat

> 📖 **Detailed Guides**: [Key Cycles](docs/key-cycles.md) | [Arrangement System](docs/arrangement.md)

### Overlay System

**Overlays** add mathematical variations to patterns:

- `1/1`: Root overlay
- `2/1`: Every 2nd cycle
- `3/1`: Every 3rd cycle
- `1/4`, `2/4`, `3/4`, `4/4`: First, second, third, fourth of every 4 cycles

Complex overlay keys support width, start delays, and stacking behaviors.

> 📖 **Detailed Guide**: [Overlay System Documentation](docs/overlay-key.md)

### Actions

**Actions** manipulate playback cursors:
- **Line Reset** (`ss`): Reset cursor to first beat
- **Line Bounce** (`sb`): Bounce between first beat and action
- **Line Skip** (`sk`): Skip current beat
- **Line Reverse** (`sr`): Reverse playback direction
- **Line Delay** (`sz`): Repeat current beat

Some actions have "All" variants (`sS`, `sB`, `sK`) that affect all playback cursors.

> 📖 **Detailed Guide**: [Actions Documentation](docs/actions.md)

## Key Bindings

### Navigation
- `hjkl`: Move cursor (left/down/up/right)
- `<>`: Jump to line start/end
- `bf`/`bl`: Jump to first/last line

### Playback
- `Space`: Play/stop arrangement once
- `Alt+Space`: Loop arrangement
- `'+Space`: Loop current overlay
- `Ctrl+@`: Loop current part

### Pattern Creation
- `1-9`: Add/remove notes every N beats
- `f`: Add single note
- `d`: Remove note
- `c`: Clear line from cursor to end

### Note Modification
- `A`/`a`: Increase/decrease accent
- `G`/`g`: Increase/decrease gate
- `R`/`r`: Increase/decrease ratchet
- `W`/`w`: Increase/decrease wait
- `E`/`e`: Long gate increase/decrease

### Arrangement
- `Ctrl+a`: Toggle arrangement view
- `Ctrl+]`: New section after current
- `Ctrl+p`: New section before current
- `]s`/`[s`: Next/previous section

### Advanced Features
- `o`: Toggle chord mode
- `v`: Visual selection mode
- `y`/`p`: Yank/paste
- `m`/`M`: Mute/solo line
- `u`/`U`: Undo/redo

### Input Modes
- `Ctrl+b`: Beat input controls
- `Ctrl+e`: Accent input controls
- `Ctrl+o`: Overlay key controls
- `Ctrl+d`: MIDI setup controls
- `Ctrl+t`: Tempo controls

> 📖 **Complete Reference**: [Key Mappings Documentation](docs/key-mappings.md)

## Advanced Features

### Chord Mode (`o`)

Create and manipulate chords with dedicated key bindings:
- `tM`/`tm`: Major/minor triads
- `7M`/`7m`: Major/minor sevenths
- `]i`/`[i`: Invert chords
- `]d`/`[d`: Double notes
- `]p`/`[p`: Arpeggio patterns

### Visual Mode (`v`)

Select and manipulate regions:
- `v`: Enter visual mode
- `hjkl`: Expand selection
- `y`: Yank selection
- `p`: Paste at cursor

### Pattern Rotation

Rotate patterns in any direction:
- `HJKL`: Rotate left/down/up/right
- Works on current line or column

## File Operations

- `Ctrl+n`: New sequence
- `Ctrl+s`: Save sequence
- `br`: Reload file (lose unsaved changes)
- `q`: Quit

## Configuration

SEQ uses Lua for configuration. Configuration files are located in the `/config` directory:

- `config/init.lua`: Main configuration file
- Templates and instruments can be added via Lua scripts

## Development

### Building

```bash

# All commands must use -tags lua54 
go build -tags lua54 -o seq
go test -tags lua54 ./...
go fmt ./...

# Consider exporting a GOFLAGS variable to avoid including `-tags lua54` for every command
export GOFLAGS="-tags=lua54"
```

### Project Structure

```
seq/
├── main.go              # Main application entry
├── ui.go               # Bubble Tea UI implementation
├── internal/           # Internal packages
│   ├── arrangement/    # Arrangement logic and views
│   ├── overlays/       # Overlay system and chord handling
│   ├── overlaykey/     # Overlay key parsing and UI
│   ├── theory/         # Music theory utilities
│   ├── seqmidi/        # MIDI connection handling
│   └── themes/         # Color themes and display
└── docs/              # Documentation
```

## Documentation

### Complete Documentation Index

- **[Beat Creation](docs/beat-creation.md)** - Pattern mode and rapid beat creation techniques
- **[Note Alteration](docs/note-alteration.md)** - Accent, gate, ratchet, and wait controls
- **[Key Cycles](docs/key-cycles.md)** - Understanding sequence iteration and timing
- **[Arrangement System](docs/arrangement.md)** - Sections, parts, groups, and song structure
- **[Overlay System](docs/overlay-key.md)** - Mathematical variations and overlay keys
- **[Actions](docs/actions.md)** - Playback cursor manipulation and sequencer actions
- **[Key Mappings](docs/key-mappings.md)** - Complete keyboard shortcut reference

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes following Go conventions
4. Add tests for new features
5. Submit a pull request

### Code Style

- Use Go 1.24+ features
- Follow standard Go naming conventions
- Use `stretchr/testify/assert` for tests
- Run `go fmt` before committing
- Include `-tags lua54` in all go commands

## License

[License information to be added]

## Support

For issues and feature requests, please use the GitHub issue tracker.
