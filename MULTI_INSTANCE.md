# TUIWrite Multi-Instance Usage Guide

**Last Updated**: November 13, 2025

---

## Overview

TUIWrite supports opening multiple files by running multiple instances. When running inside `tmux` or `screen`, you can open files in split panes for side-by-side editing.

---

## Commands

### Opening Files

#### `:e <filename>` (or `:edit`, `:open`)
Opens a file in the **current instance**, replacing the current document.

```
:e chapter2.md
:edit notes.txt
:open script.fountain
```

- If the file doesn't exist, it will be created
- Unsaved changes are **not** automatically saved (you'll lose them)
- File tree is closed automatically when opening

**Use case**: Quick switching between files in the same window

---

#### `:new <filename>` (or `:vnew`)
Opens a file in a **new instance**.

**Behavior depends on environment:**

**In tmux** (recommended):
- Opens file in a **side-by-side split** (vertical split)
- New pane appears to the right
- Easy to switch with `Ctrl+B`, then arrow keys

```
:new chapter2.md    # Opens chapter2.md in right pane
:vnew notes.txt     # Same as :new
```

**In screen**:
- Opens file in a new split
- Switch panes with `Ctrl+A`, then Tab

**Outside tmux/screen**:
- Shows message: "Not in tmux/screen. Use :e to open in current instance"
- Suggests starting tmux first

---

#### `:split <filename>`
Opens a file in a **new instance** with top/bottom split.

**In tmux**:
- Opens file in a **horizontal split** (top/bottom)
- New pane appears below current pane

```
:split reference.md    # Opens reference.md below current pane
```

**Outside tmux**:
- Same behavior as `:new` (shows message)

---

### File Tree Integration

When the **file tree is visible** (press `F1`):

1. Navigate with **arrow keys** or **j/k**
2. Press **Enter** on a file to open it in the **current instance**
3. Directories can be expanded/collapsed with **Enter**

To open a file from the tree in a **new instance**:
1. Note the filename
2. Press `Esc` to exit file tree
3. Use `:new <filename>`

---

### Getting Help

#### `:help` (or `:h`)
Shows available commands and multiplexer status.

Output examples:
- **In tmux**: "Commands: :e <file> (this window) | :new <file> (side split) | :split <file> (top/bottom)"
- **In screen**: "Commands: :e <file> (this window) | :new <file> (new split)"
- **Standalone**: "Commands: :e <file> (open in current instance) | Use tmux/screen for :new"

---

## Workflow Examples

### Example 1: Writing with Reference Material (tmux)

```bash
# Start tmux
tmux

# Open your main document
./tuiwrite novel.md

# Open reference material in a side pane
:new outline.md

# Now you have novel.md on the left, outline.md on the right
# Switch between them with Ctrl+B, then arrow keys
```

### Example 2: Screenplay with Notes (tmux vertical split)

```bash
tmux
./tuiwrite script.fountain

# Open notes below the script
:split notes.md

# script.fountain is on top, notes.md on bottom
```

### Example 3: Quick File Switching (single instance)

```bash
./tuiwrite chapter1.md

# Switch to chapter 2 without opening new instance
:e chapter2.md

# Go back to chapter 1
:e chapter1.md
```

### Example 4: Using File Tree

```bash
./tuiwrite

# Press F1 to open file tree
# Navigate to desired file
# Press Enter to open it
# Press Esc to return to Read mode
```

---

## Tmux Basics (if new to tmux)

### Starting tmux
```bash
tmux                    # Start new session
tmux new -s writing     # Start named session
```

### Switching between panes
- `Ctrl+B`, then **arrow keys** - Switch to pane in that direction
- `Ctrl+B`, then `o` - Switch to next pane
- `Ctrl+B`, then `q` - Show pane numbers

### Resizing panes
- `Ctrl+B`, then `Ctrl+arrow` - Resize current pane

### Closing panes
- In TUIWrite: `:q` or `Ctrl+Q`
- In terminal: `Ctrl+D` or `exit`

### Detaching from tmux
- `Ctrl+B`, then `d` - Detach (session keeps running)
- `tmux attach` - Reattach to session

---

## Screen Basics (alternative to tmux)

### Starting screen
```bash
screen                  # Start new session
```

### Switching between regions
- `Ctrl+A`, then `Tab` - Switch to next region

### Creating splits
- TUIWrite handles this with `:new`

### Closing regions
- In TUIWrite: `:q` or `Ctrl+Q`

### Detaching from screen
- `Ctrl+A`, then `d` - Detach
- `screen -r` - Reattach

---

## Tips & Tricks

### 1. **Use tmux for best experience**
tmux is more actively developed and has better split support than screen.

```bash
# Install tmux
sudo apt install tmux        # Ubuntu/Debian
sudo dnf install tmux        # Fedora
brew install tmux            # macOS
```

### 2. **Create a writing session script**

Create `~/bin/writing-session`:
```bash
#!/bin/bash
tmux new-session -d -s writing "cd ~/Documents/writing && tuiwrite novel.md"
tmux split-window -h -t writing "cd ~/Documents/writing && tuiwrite outline.md"
tmux attach -t writing
```

Make it executable:
```bash
chmod +x ~/bin/writing-session
```

Run it:
```bash
writing-session
```

### 3. **Quick file tree + open workflow**
1. Press `F1` to open file tree
2. Navigate to file
3. Press `Enter` to open in current instance
4. File tree auto-closes

### 4. **Keep a scratch buffer open**
In one tmux pane, keep an untitled document for quick notes:
```bash
tmux
./tuiwrite              # Untitled document
# In another pane:
:split main-document.md
```

### 5. **Don't forget Ctrl+Q**
Since `Ctrl+C` is now copy, use **`Ctrl+Q`** to quit.

---

## Limitations & Future Enhancements

### Current Limitations
- `:e` doesn't prompt to save unsaved changes (will be added)
- No way to list all open instances
- Can't share undo/redo between instances
- Each instance loads its own spell-check dictionary

### Planned Enhancements (v0.3)
- `:e!` to force-open without saving
- `:bd` (buffer delete) as alias for `:q`
- Recent files list (`:recent`)
- Session management (save/restore tmux layouts)

---

## Troubleshooting

### "Not in tmux/screen" message
**Problem**: `:new` command shows this message.

**Solution**: Start tmux or screen:
```bash
tmux
./tuiwrite yourfile.md
```

### File opens but I can't see it
**Problem**: New tmux pane opened but isn't visible.

**Solution**: Use `Ctrl+B`, then arrow keys to find the pane.

### File tree won't open files
**Problem**: Pressing Enter on files does nothing.

**Solution**: Make sure you're not in command mode. Press `Esc` first, then try again.

### Can't quit with Ctrl+C
**Problem**: `Ctrl+C` doesn't quit anymore.

**Solution**: This is intentional - `Ctrl+C` is now copy. Use **`Ctrl+Q`** to quit.

---

## Command Reference Summary

| Command | Description | Example |
|---------|-------------|---------|
| `:e <file>` | Open file in current instance | `:e chapter2.md` |
| `:edit <file>` | Alias for `:e` | `:edit notes.txt` |
| `:open <file>` | Alias for `:e` | `:open script.fountain` |
| `:new <file>` | Open file in new side-by-side split (tmux) | `:new outline.md` |
| `:vnew <file>` | Alias for `:new` | `:vnew reference.md` |
| `:split <file>` | Open file in new top/bottom split (tmux) | `:split notes.md` |
| `:help` | Show command help and multiplexer status | `:help` |
| `:h` | Alias for `:help` | `:h` |
| `:w` | Save current file | `:w` |
| `:q` | Quit current instance | `:q` |
| `:wq` | Save and quit | `:wq` |
| `F1` | Toggle file tree | Press `F1` |
| `Ctrl+Q` | Quit | Press together |

---

**End of Multi-Instance Guide**
