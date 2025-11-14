# TUIWrite Quick Reference Card

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           TUIWRITE QUICK REFERENCE                       │
│                         Modal TUI Text Editor                            │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ MODES                                                                    │
├─────────────────────────────────────────────────────────────────────────┤
│ Read Mode    Navigation and commands (default)                          │
│ Edit Mode    Full text editing                                          │
│ Command Mode Execute commands (start with :)                            │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ MODE SWITCHING                                                           │
├─────────────────────────────────────────────────────────────────────────┤
│ INSERT       Enter Edit mode                                            │
│ ESC          Return to Read mode                                        │
│ :            Enter Command mode (from Read mode)                        │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ NAVIGATION (Both Modes)                                                 │
├─────────────────────────────────────────────────────────────────────────┤
│ ↑↓←→         Navigate by visual wrapped lines                           │
│ h j k l      Vim-style navigation                                       │
│ HOME / END   Start/end of line                                          │
│ PgUp / PgDn  Page up/down                                               │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ EDITING (Edit Mode)                                                     │
├─────────────────────────────────────────────────────────────────────────┤
│ Type         Insert characters at cursor                                │
│ ENTER        New line                                                   │
│ BACKSPACE    Delete character before cursor                             │
│ DELETE       Delete character at cursor                                 │
│ TAB          Insert 4 spaces                                            │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ COPY/PASTE (Edit Mode)                                                  │
├─────────────────────────────────────────────────────────────────────────┤
│ SHIFT+↑↓←→   Select text (keyboard)                                     │
│ Click+Drag   Select text (mouse)                                        │
│ Ctrl+C       Copy selected text to clipboard                            │
│ Ctrl+V       Paste from clipboard                                       │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ GLOBAL SHORTCUTS                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│ Ctrl+S       Save file                                                  │
│ Ctrl+Q       Quit                                                       │
│ F1           Toggle file tree                                           │
│ F11          Toggle zen mode (fullscreen)                               │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ FILE COMMANDS                                                            │
├─────────────────────────────────────────────────────────────────────────┤
│ :w           Save current file                                          │
│ :write       Same as :w                                                 │
│ :q           Quit current instance                                      │
│ :quit        Same as :q                                                 │
│ :wq          Save and quit                                              │
│                                                                          │
│ :e <file>    Open file in current instance                              │
│ :edit <file> Same as :e                                                 │
│ :open <file> Same as :e                                                 │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ MULTI-INSTANCE COMMANDS (tmux/screen)                                   │
├─────────────────────────────────────────────────────────────────────────┤
│ :new <file>  Open file in new side-by-side split                        │
│ :vnew <file> Same as :new                                               │
│ :split <file> Open file in new top/bottom split                         │
│                                                                          │
│ Note: Requires tmux or screen. Otherwise shows helpful message.         │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ SPELL-CHECK COMMANDS                                                     │
├─────────────────────────────────────────────────────────────────────────┤
│ :spell       Toggle spell-check on/off                                  │
│ :spellcheck  Same as :spell                                             │
│ :spell -uk   Enable UK English spell-check                              │
│ :spell -us   Enable US English spell-check                              │
│                                                                          │
│ Supported: uk, gb, us, ca, au, es, fr, de, it, pt                       │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ HELP COMMAND                                                             │
├─────────────────────────────────────────────────────────────────────────┤
│ :help        Show command help and multiplexer status                   │
│ :h           Same as :help                                              │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ FILE TREE (F1)                                                           │
├─────────────────────────────────────────────────────────────────────────┤
│ F1           Toggle file tree visibility                                │
│ ↑↓ or j/k    Navigate files/folders                                     │
│ ENTER        Open file (current instance) or toggle folder              │
│ ESC          Close file tree / return focus to editor                   │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ COMMON WORKFLOWS                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│ Quick Edit                                                               │
│   1. Press INSERT to enter Edit mode                                    │
│   2. Make changes                                                        │
│   3. Press Ctrl+S to save                                               │
│   4. Press ESC to return to Read mode                                   │
│                                                                          │
│ Browse and Open File                                                     │
│   1. Press F1 to show file tree                                         │
│   2. Use ↑↓ to find file                                                │
│   3. Press ENTER to open                                                │
│                                                                          │
│ Side-by-Side Editing (tmux)                                             │
│   1. Start tmux                                                          │
│   2. Open main file: ./tuiwrite novel.md                                │
│   3. Type :new outline.md                                               │
│   4. Switch panes: Ctrl+B then arrow keys                               │
│                                                                          │
│ Copy Between Documents                                                   │
│   1. Select text with SHIFT+arrows or mouse                             │
│   2. Press Ctrl+C to copy                                               │
│   3. Switch to other instance/pane                                      │
│   4. Press INSERT then Ctrl+V to paste                                  │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ TMUX BASICS (for multi-instance)                                        │
├─────────────────────────────────────────────────────────────────────────┤
│ tmux         Start tmux session                                         │
│ Ctrl+B, →    Switch to right pane                                       │
│ Ctrl+B, ←    Switch to left pane                                        │
│ Ctrl+B, ↑    Switch to upper pane                                       │
│ Ctrl+B, ↓    Switch to lower pane                                       │
│ Ctrl+B, o    Cycle through panes                                        │
│ Ctrl+B, d    Detach from session (keeps running)                        │
│ tmux attach  Re-attach to session                                       │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ STATUS BAR COLORS                                                        │
├─────────────────────────────────────────────────────────────────────────┤
│ Green        Success (file saved, spell-check enabled, etc.)            │
│ Yellow       Warning (no selection, missing file, etc.)                 │
│ Red          Error (save failed, invalid command, etc.)                 │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ TIPS                                                                     │
├─────────────────────────────────────────────────────────────────────────┤
│ • Ctrl+Q to quit (not Ctrl+C - that's copy now!)                        │
│ • Use tmux for best multi-instance experience                           │
│ • File tree auto-closes when opening files                              │
│ • Spell-check dictionaries download automatically first time            │
│ • Selection clears when you type or move without SHIFT                  │
│ • Mouse selection works if your terminal supports it                    │
└─────────────────────────────────────────────────────────────────────────┘

                           Version 0.2 - November 2025
                      See MULTI_INSTANCE.md for detailed guide
```
