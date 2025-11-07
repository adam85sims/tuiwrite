# Debug Logging

TUIWrite includes comprehensive debug logging to help diagnose issues and track program behavior.

## Log File Location

The debug log is automatically created in your platform's standard config directory:

- **Linux/Unix**: `~/.config/tuiwrite/debug.log`
- **macOS**: `~/Library/Application Support/tuiwrite/debug.log`
- **Windows**: `%APPDATA%\tuiwrite\debug.log`

## What Gets Logged

The logger tracks all significant operations:

### Session Events
- Program start/exit
- Platform information
- File being edited
- Document mode (story/script)

### File Operations
- File loading (with line count)
- File saving (with byte count)
- File creation
- Save errors

### User Actions
- Mode switches (Read ↔ Edit)
- Command execution (`:w`, `:q`, `:spellcheck`, etc.)
- Manual saves (Ctrl+S)
- Program quit

### System Events
- Window resizing
- Dictionary downloads
- Spell-check language changes
- Cache invalidations

### Errors
- File I/O errors
- Dictionary download failures
- Save failures
- Unknown commands

## Log Levels

Logs are categorized by severity:

- **SESSION_START/END**: Program lifecycle events
- **EVENT**: User actions and mode changes
- **INFO**: Normal operations (file loads, saves)
- **DEBUG**: Detailed information (window size, cache operations)
- **WARNING**: Non-critical issues (missing dictionaries, unknown commands)
- **ERROR**: Failures (file save errors, network issues)

## Example Log Output

```
[2025-11-07 09:50:18.391] SESSION_START: TUIWrite started
[2025-11-07 09:50:18.391] INFO: Log file: /home/adam/.config/tuiwrite/debug.log
[2025-11-07 09:50:18.391] INFO: Platform: linux/amd64
[2025-11-07 09:50:18.391] INFO: Starting TUIWrite with file: README.md
[2025-11-07 09:50:18.391] INFO: Document mode: story
[2025-11-07 09:50:18.391] DEBUG: Loading file: README.md
[2025-11-07 09:50:18.392] DEBUG: Loaded 150 lines from README.md
[2025-11-07 09:50:18.392] INFO: File loaded: 150 lines
[2025-11-07 09:50:18.392] INFO: Starting Bubble Tea program
[2025-11-07 09:50:18.392] INFO: Model initialized
[2025-11-07 09:50:18.392] INFO: File: README.md, Mode: story
[2025-11-07 09:50:18.392] DEBUG: Window resized: 103x24
[2025-11-07 09:50:23.680] EVENT: MODE_CHANGE: Switched to EDIT mode
[2025-11-07 09:50:54.735] EVENT: SAVE: Manual save triggered
[2025-11-07 09:50:54.735] DEBUG: Saving file: README.md (150 lines)
[2025-11-07 09:50:54.735] INFO: File saved successfully: README.md (5247 bytes)
[2025-11-07 09:50:54.735] INFO: File saved successfully
[2025-11-07 09:51:01.114] INFO: User requested quit
[2025-11-07 09:51:01.124] INFO: Program exited normally
[2025-11-07 09:51:01.124] SESSION_END: TUIWrite exiting
```

## Using the Logs for Debugging

### Finding Issues
View the most recent log entries:
```bash
tail -50 ~/.config/tuiwrite/debug.log
```

### Tracking a Specific Session
Each session is bookended with `SESSION_START` and `SESSION_END`:
```bash
grep -A 100 "SESSION_START" ~/.config/tuiwrite/debug.log | tail -100
```

### Finding Errors
Filter for errors and warnings:
```bash
grep -E "ERROR|WARNING" ~/.config/tuiwrite/debug.log
```

### Performance Analysis
Timestamps show exact timing of operations:
```bash
grep "SAVE" ~/.config/tuiwrite/debug.log
```

### Clearing the Log
If the log gets too large, you can safely delete it:
```bash
rm ~/.config/tuiwrite/debug.log
```

A new log will be created automatically on next program start.

## Log Rotation

Currently, the log appends indefinitely. For very long-term use, you may want to periodically clear or rotate the log file. Future versions may include automatic log rotation.

## Privacy Note

The debug log contains:
- File paths you've edited
- Commands you've executed
- Timing information

It does NOT contain:
- File contents
- Passwords or sensitive data
- Network traffic beyond connectivity checks

The log is stored locally and never transmitted anywhere.
