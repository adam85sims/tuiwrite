# tuiwrite - First Draft

A modal TUI text editor for prose and screenplay writing, built with Go and Bubble Tea.

**Performance**: Optimized for large documents with lazy word wrapping (handles 10,000+ line files smoothly)

## Features Implemented (v0.2)

### Core Functionality
- ✅ Modal interface (Read Mode / Edit Mode)
- ✅ File loading and creation
- ✅ Manual save (Ctrl+S)
- ✅ Auto-save every 30 minutes
- ✅ Basic text editing (insert, delete, backspace)
- ✅ Line navigation with arrow keys
- ✅ Page navigation (PgUp/PgDn)
- ✅ Document start/end navigation (Ctrl+Home/End, g/G)
- ✅ Status bar with mode indicator, filename, and position
- ✅ Modified indicator [+]
- ✅ Colored status messages
- ✅ **Lazy word wrapping** (3,192x faster for large files!)
- ✅ Visible cursor in both Read and Edit modes
- ✅ Command mode (`:` in Read Mode)
- ✅ On-demand spell-checking with dictionary downloads

### Performance
- **Lazy wrapping**: Only processes visible lines (not entire document)
- **Per-line caching**: Wrapped results cached and reused
- **Smart invalidation**: Only edited lines re-wrapped
- **Benchmarks**: 440ns rendering vs 1.4ms (old approach)
- **Memory efficient**: 1KB per render vs 1.7MB (99.94% reduction)
- **Tested**: Handles 10,000 line documents with zero lag

### Debug & Diagnostics
- **Comprehensive logging**: All operations logged to `debug.log`
- **Cross-platform**: Logs stored in standard config directories
  - Linux: `~/.config/tuiwrite/debug.log`
  - macOS: `~/Library/Application Support/tuiwrite/debug.log`
  - Windows: `%APPDATA%\tuiwrite\debug.log`
- **Timestamped**: Millisecond-precision timestamps for debugging
- **Event tracking**: Mode changes, commands, file ops, errors
- See [LOGGING.md](LOGGING.md) for details

## Usage

### Building
```bash
go build -o tuiwrite main.go
```

### Running
```bash
# Create/open a story document
./tuiwrite mynovel.md

# Create/open a screenplay document
./tuiwrite myscript.fountain -mode script
```

## Keybindings

### Global (Both Modes)
- `Ctrl+S` - Save file
- `Ctrl+C` / `Ctrl+Q` - Quit
- `Insert` - Toggle to Edit Mode (from Read) or Enter Edit Mode
- `Esc` - Switch to Read Mode (from Edit)

### Read Mode (Navigation Only)
- `↑↓←→` or `hjkl` - Navigate cursor (visible but read-only)
- `PgUp/PgDn` - Scroll by page
- `Ctrl+Home` / `g` - Jump to document start
- `Ctrl+End` / `G` - Jump to document end
- `Home` - Start of line
- `End` - End of line
- `:` - Enter command mode

### Edit Mode (Full Editing)
- `↑↓←→` - Navigate cursor
- `Enter` - New line
- `Backspace` - Delete character before cursor
- `Delete` - Delete character at cursor
- `Home` - Start of line
- `End` - End of line
- Any printable character - Insert at cursor

## Command Mode

Press `:` in Read Mode to enter command mode. Available commands:

### Spell-Checking Commands
- `:spellcheck` or `:spell` - Toggle spell checking on/off
- `:spellcheck -uk` or `:spell -uk` - Enable UK English spell-check (default)
- `:spellcheck -us` - Enable US English spell-check
- `:spellcheck -ca` - Enable Canadian English
- `:spellcheck -au` - Enable Australian English
- `:spellcheck -es` - Enable Spanish
- `:spellcheck -fr` - Enable French
- `:spellcheck -de` - Enable German
- `:spellcheck -it` - Enable Italian
- `:spellcheck -pt` - Enable Portuguese

**Note:** Dictionaries are downloaded automatically on first use and cached in `~/.config/tuiwrite/dictionaries/`

### File Commands
- `:w` or `:write` - Save file
- `:q` or `:quit` - Quit application
- `:wq` - Save and quit

Press `Esc` to exit command mode.

## Status Bar

The status bar shows:
- Current mode (READ or EDIT)
- Filename and modified indicator [+]
- Current line and column position
- Status messages (saves, errors, downloads)
- Command buffer (when in command mode)
- Quick help text

## Spell-Checking

TUIWrite includes built-in spell-checking powered by Hunspell dictionaries:

- **On-demand downloads**: Dictionaries are automatically downloaded when you first select a language
- **Multiple languages**: Support for 10 languages (see command mode above)
- **Lightweight**: Only downloads dictionaries you actually use
- **Cached locally**: Dictionaries are stored locally for offline use (see paths below)
- **Dictionary source**: [https://github.com/adam85sims/tuiwritedics](https://github.com/adam85sims/tuiwritedics)
- **Internet required**: First-time dictionary downloads require internet connectivity

### Dictionary Storage Locations

Dictionaries are cached in platform-specific locations:

- **Linux/Unix**: `~/.config/tuiwrite/dictionaries/`
- **macOS**: `~/Library/Application Support/tuiwrite/dictionaries/`
- **Windows**: `%APPDATA%\tuiwrite\dictionaries\`

Once downloaded, dictionaries work offline. If you're not connected to the internet when trying to download a new dictionary, you'll receive an error message asking you to connect and try again.

Visual spell-check indicators (underlining misspelled words) will be added in a future update.

## Features Not Yet Implemented

The following features from the design document are planned for future versions:
- F1-F4 function keys (file tree, statistics, search, menu)
- Visual spell-check highlighting
- Formatting commands (#bold:, #italics:, etc.)
- Structural commands (#title:, #chapter:, #break:, etc.)
- File tree sidebar
- Statistics panel
- Search/find functionality
- Export functionality
- Screenplay-specific formatting

## Notes

This is a first draft implementation focusing on the core modal editing experience and basic file operations. The application starts in Read Mode by default to prevent accidental edits, as specified in the design document.
[Logging Test]