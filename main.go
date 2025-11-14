package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the model
func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tickAutoSave(),
		tickCursorBlink(),
	)
}

// tickAutoSave returns a command that sends autoSaveMsg every 30 minutes
func tickAutoSave() tea.Cmd {
	return tea.Tick(30*time.Minute, func(t time.Time) tea.Msg {
		return autoSaveMsg(t)
	})
}

// tickCursorBlink returns a command that sends cursorBlinkMsg every 500ms
func tickCursorBlink() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return cursorBlinkMsg(t)
	})
}

// Update handles messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Reset cursor to visible when user types
		m.cursorVisible = true
		m.lastCursorBlink = time.Now()
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Recalculate word wrapping when window size changes
		m.rewrapLines()
		return m, nil

	case autoSaveMsg:
		if m.modified {
			err := m.saveFile()
			if err == nil {
				m.setStatus("Auto-saved", "green")
			}
		}
		return m, tickAutoSave()

	case cursorBlinkMsg:
		// Toggle cursor visibility
		m.cursorVisible = !m.cursorVisible
		m.lastCursorBlink = time.Now()
		return m, tickCursorBlink()

	case tea.MouseMsg:
		return m.handleMouse(msg)
	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle command mode separately
	if m.commandMode {
		return m.handleCommandMode(msg)
	}

	// Global keybindings (work in both modes)
	switch msg.String() {
	case "ctrl+q":
		return m, tea.Quit

	case "ctrl+s":
		err := m.saveFile()
		if err != nil {
			m.setStatus("Error saving: "+err.Error(), "red")
		} else {
			m.setStatus("Saved "+m.filename, "green")
		}
		return m, nil

	case "f1":
		// Toggle file tree sidebar
		return m.toggleFileTree()

	case "insert":
		// Toggle between Read and Edit mode (only if file tree not focused)
		if m.fileTreeFocused {
			return m, nil
		}
		if m.mode == ReadMode {
			m.mode = EditMode
			m.setStatus("-- EDIT MODE --", "green")
		} else {
			m.mode = ReadMode
			m.setStatus("-- READ MODE --", "green")
		}
		return m, nil

	case "esc":
		if m.mode == EditMode {
			m.mode = ReadMode
			m.setStatus("-- READ MODE --", "green")
		}
		return m, nil

	case ":":
		// Enter command mode (only in Read mode and file tree not focused)
		if m.mode == ReadMode && !m.fileTreeFocused {
			m.commandMode = true
			m.commandBuffer = ":"
			return m, nil
		}
		// If in Edit mode, fall through to normal edit handling
	}

	// If file tree is focused, handle file tree navigation
	if m.fileTreeFocused {
		return m.handleFileTreeNavigation(msg.String())
	}

	// Mode-specific keybindings
	if m.mode == ReadMode {
		return m.handleReadMode(msg)
	} else {
		return m.handleEditMode(msg)
	}
}

// handleCommandMode processes keys when in command mode
func (m model) handleCommandMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Exit command mode
		m.commandMode = false
		m.commandBuffer = ""
		return m, nil

	case "enter":
		// Execute command
		cmd := m.commandBuffer
		m.commandMode = false
		m.commandBuffer = ""
		return m.executeCommand(cmd)

	case "backspace":
		// Delete character from command buffer
		if len(m.commandBuffer) > 1 {
			m.commandBuffer = m.commandBuffer[:len(m.commandBuffer)-1]
		} else {
			// If only ":" left, exit command mode
			m.commandMode = false
			m.commandBuffer = ""
		}
		return m, nil

	default:
		// Add character to command buffer
		if len(msg.Runes) == 1 {
			m.commandBuffer += string(msg.Runes[0])
		}
		return m, nil
	}
}

// executeCommand executes a command entered in command mode
func (m model) executeCommand(cmd string) (tea.Model, tea.Cmd) {
	// Remove the leading ":"
	if len(cmd) > 1 {
		cmd = cmd[1:]
	} else {
		return m, nil
	}

	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return m, nil
	}

	switch parts[0] {
	case "spellcheck", "spell":
		if len(parts) == 1 {
			// Toggle spell checking
			m.spellChecker.toggle()
			if m.spellChecker.enabled {
				m.setStatus("Spell checking enabled ("+m.spellChecker.language+")", "green")
			} else {
				m.setStatus("Spell checking disabled", "yellow")
			}
		} else if len(parts) == 2 {
			// Set language
			lang := strings.TrimPrefix(parts[1], "-")
			lang = strings.ToLower(lang)

			// Check if dictionary needs downloading
			if !m.spellChecker.hasDictionary(lang) {
				m.setStatus("Downloading "+strings.ToUpper(lang)+" dictionary...", "yellow")
			}

			// Set the language (will download if needed)
			err := m.spellChecker.setLanguage(lang)
			if err != nil {
				m.setStatus("Failed to load dictionary: "+err.Error(), "red")
			} else {
				m.setStatus("Spell-check enabled ("+strings.ToUpper(lang)+")", "green")
			}
		}
		return m, nil

	case "q", "quit":
		return m, tea.Quit

	case "w", "write":
		err := m.saveFile()
		if err != nil {
			m.setStatus("Error saving: "+err.Error(), "red")
		} else {
			m.setStatus("Saved "+m.filename, "green")
		}
		return m, nil

	case "wq":
		err := m.saveFile()
		if err != nil {
			m.setStatus("Error saving: "+err.Error(), "red")
			return m, nil
		}
		return m, tea.Quit

	case "e", "edit", "open":
		// Open file in current instance
		if len(parts) < 2 {
			m.setStatus("Usage: :e <filename>", "yellow")
			return m, nil
		}
		return m.openFileInCurrentInstance(parts[1])

	case "new", "vnew", "split":
		// Open file in new instance (tmux/screen split if available)
		if len(parts) < 2 {
			m.setStatus("Usage: :new <filename>", "yellow")
			return m, nil
		}
		return m.openFileInNewInstance(parts[1], parts[0])

	case "help", "h":
		// Show command help
		return m.showMultiplexerHelp()

	default:
		m.setStatus("Unknown command: "+parts[0], "red")
		return m, nil
	}
}

// handleMouse processes mouse events for text selection
func (m model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Only handle mouse in edit mode
	if m.mode != EditMode {
		return m, nil
	}

	// Skip if file tree is focused
	if m.fileTreeFocused {
		return m, nil
	}

	switch msg.Button {
	case tea.MouseButtonLeft:
		if msg.Action == tea.MouseActionPress {
			// Mouse down - start selection
			// Convert screen coordinates to document coordinates
			// msg.Y is screen row, need to account for status bars and viewport offset
			// First 2 lines are status, so document starts at line 2
			screenY := msg.Y - 2 // Account for status bars at top
			if screenY < 0 {
				return m, nil // Click was in status area
			}

			// Convert wrapped line position to source line
			wrappedIdx := m.offsetY + screenY
			sourceY, offsetInLine := m.getSourceLineFromWrappedIndex(wrappedIdx)

			// Calculate X position (accounting for horizontal scroll)
			clickX := msg.X + m.offsetX + offsetInLine

			// Clamp to line bounds
			if sourceY >= 0 && sourceY < len(m.lines) {
				lineLen := len(m.lines[sourceY])
				if clickX > lineLen {
					clickX = lineLen
				}

				m.cursorY = sourceY
				m.cursorX = clickX
				m.selectionActive = true
				m.selectionStartX = clickX
				m.selectionStartY = sourceY
			}
		} else if msg.Action == tea.MouseActionMotion && m.selectionActive {
			// Mouse drag - extend selection
			screenY := msg.Y - 2
			if screenY < 0 {
				return m, nil
			}

			wrappedIdx := m.offsetY + screenY
			sourceY, offsetInLine := m.getSourceLineFromWrappedIndex(wrappedIdx)
			clickX := msg.X + m.offsetX + offsetInLine

			if sourceY >= 0 && sourceY < len(m.lines) {
				lineLen := len(m.lines[sourceY])
				if clickX > lineLen {
					clickX = lineLen
				}

				m.cursorY = sourceY
				m.cursorX = clickX
			}
		}
	}

	return m, nil
}

// getUntitledFilename generates a unique untitled document filename
func getUntitledFilename() string {
	// Find the next available untitled-document-N
	counter := 1
	for {
		filename := fmt.Sprintf("untitled-document-%d", counter)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return filename
		}
		counter++
	}
}

func main() {
	// Initialize logger
	if err := initLogger(); err != nil {
		fmt.Printf("Warning: Failed to initialize logger: %v\n", err)
	}
	defer CloseLogger()

	// Parse command line arguments
	modeFlag := flag.String("mode", "story", "Document mode: story or script")
	flag.Parse()

	args := flag.Args()

	var filename string
	var lines []string
	var isUntitled bool

	if len(args) == 0 {
		// No filename provided - create untitled document
		filename = getUntitledFilename()
		lines = []string{""}
		isUntitled = true
		LogInfo("Starting with untitled document: " + filename)
	} else {
		// Filename provided
		filename = args[0]
		isUntitled = false

		// Load or create file
		var err error
		lines, err = loadFile(filename)
		if err != nil {
			fmt.Printf("Error loading file: %v\n", err)
			os.Exit(1)
		}
		LogInfo("Loaded file: " + filename)
	}

	// Determine document mode
	docMode := StoryMode
	if *modeFlag == "script" {
		docMode = ScriptMode
	}

	// Initialize model
	m := model{
		filename:          filename,
		docMode:           docMode,
		lines:             lines,
		cursorX:           0,
		cursorY:           0,
		offsetY:           0,
		offsetX:           0,
		mode:              ReadMode,    // Start in read mode
		saved:             !isUntitled, // Untitled docs start as unsaved
		modified:          isUntitled,  // Untitled docs start as modified
		lastSave:          time.Now(),
		wrapCache:         make(map[int][]wrappedLine),
		wrapWidth:         0,
		fontSize:          DefaultFontSize, // 100%
		fontSizeDirection: "",
		lastFontSizeTime:  time.Now(),
		zenMode:           false,
		cursorVisible:     true, // Start with cursor visible
		lastCursorBlink:   time.Now(),
		spellChecker:      newSpellChecker("uk"), // Default to UK English
		commandMode:       false,
		commandBuffer:     "",
		fileTreeVisible:   false,
		fileTreeFocused:   false,
		fileTreeCursor:    0,
		fileTreeOffset:    0,
	}

	// Initialize file tree
	if err := m.initFileTree(); err != nil {
		LogWarningf("Failed to initialize file tree: %v", err)
	}

	// Run the program
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
