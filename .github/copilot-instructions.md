# TUIWrite AI Coding Assistant Guide

## Project Overview
TUIWrite is a modal TUI text editor for prose/screenplay writing built with Go and Bubble Tea. Performance-optimized for large documents (10,000+ lines) using lazy word wrapping with per-line caching.

## Core Architecture

### State Management: The `model` struct
- **Single source of truth** in `types.go` - all app state lives here
- **Cursor coordinates**: `cursorX/cursorY` are *source line* positions, `offsetY` is in *wrapped line* coordinates
- **Update pattern**: All state changes go through `Update(msg) (tea.Model, tea.Cmd)` - return new state + optional commands

### Critical: Lazy Wrapping System (`wrap.go`)
**Performance bottleneck solved**: Only wraps visible lines (~440ns vs 1.4ms for entire document)

Key invariants:
- Cache stored in `wrapCache map[int][]wrappedLine` - keyed by source line index
- **After ANY edit**: Must call `m.invalidateWrapCache(lineIdx)` for single-line edits
- **After insert/delete**: Call `m.invalidateWrapCacheFrom(lineIdx)` to invalidate from point onward
- **After resize**: `rewrapLines()` clears all cache if width changed
- **Navigation**: Use `getWrappedLineIndexForCursor()` and `moveToWrappedLine(idx)` for up/down - never manipulate cursorY directly for visual navigation

Pattern for edits:
```go
m.modified = true
m.invalidateWrapCache(m.cursorY)  // or invalidateWrapCacheFrom for structural changes
m.adjustViewport()
```

### Receiver Patterns
- `func (m model)` for **reads** - Init, Update, View, handlers that return new model
- `func (m *model)` for **mutations** - setStatus, adjustViewport, invalidateWrapCache, rewrapLines
- This is intentional - Bubble Tea pattern returns new state, internal mutations use pointers

## File Structure
- `types.go` - All type definitions, constants, and the `model` struct
- `main.go` - Entry point, Init/Update loop, command mode handling
- `wrap.go` - Lazy wrapping engine (performance-critical)
- `editor.go` - Edit mode key handling
- `navigation.go` - Read mode navigation 
- `view.go` - UI rendering with lipgloss
- `file.go` - File I/O
- `filetree.go` - Sidebar (F1)
- `spellcheck.go` - On-demand dictionary downloads
- `colors.go` - Catppuccin Mocha theme
- `logger.go` - Debug logging system

## Key Conventions

### Color Rendering
**CRITICAL**: Always use `ColorToHex()` with lipgloss, never `ColorToANSI()`
```go
// CORRECT
lipgloss.NewStyle().Foreground(lipgloss.Color(ColorToHex(Text)))

// WRONG - returns ANSI escape sequences, not hex
lipgloss.NewStyle().Foreground(lipgloss.Color(ColorToANSI(Text, true)))
```

### Debug Logging
- Platform-specific paths: `~/.config/tuiwrite/debug.log` (Linux), `~/Library/Application Support/tuiwrite/debug.log` (macOS), `%APPDATA%\tuiwrite\debug.log` (Windows)
- Use: `LogInfo()`, `LogDebug()`, `LogError()`, `LogEvent()` from `logger.go`
- Log state transitions, errors, performance-sensitive operations

### Modal Interface
- **Read Mode** (default): Navigation only, vim-style keys (hjkl, g/G)
- **Edit Mode**: Full editing
- Command mode: `:` enters command buffer in Read mode only
- File tree focus: Separate state `fileTreeFocused` - blocks other key handlers when true

## Build & Test

```bash
# Build
go build -o tuiwrite

# Run
./tuiwrite filename.md              # Story mode (default)
./tuiwrite script.fountain -mode script

# Test lazy wrapping (critical!)
go test -v -run TestLazyWrapping
go test -bench BenchmarkLazyWrapping
```

## Common Pitfalls

1. **Don't forget cache invalidation** - Every line edit needs `invalidateWrapCache(lineIdx)`
2. **Navigation by wrapped lines** - Use `moveToWrappedLine()`, not direct `cursorY++`
3. **lipgloss styling** - Must set `.Width(m.width)` for full-width backgrounds
4. **Pointer vs value receivers** - Follow existing pattern (mutations use `*model`)
5. **Bubble Tea local dependency** - Project uses `replace` directive for local `bubbletea-1.3.10/` - don't update without testing

## Adding Features

- **New mode/state**: Add to `model` in `types.go`, handle in `Update()`
- **New keybinding**: Add to `handleKeyPress()` (global) or mode-specific handlers
- **UI changes**: Edit `View()` and `renderStatusBar()` in `view.go`
- **Performance-sensitive**: Profile first - wrapping was 3,192x slower before optimization

## Documentation
See `CONTEXT.md` for recent changes, `STRUCTURE.md` for file organization, `LAZY_WRAPPING.md` for wrapping deep-dive, `designdoc.md` for feature spec.
