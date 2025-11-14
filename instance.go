package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// openFileInCurrentInstance loads a file into the current editor instance
func (m model) openFileInCurrentInstance(filename string) (tea.Model, tea.Cmd) {
	// Resolve to absolute path
	absPath, err := filepath.Abs(filename)
	if err != nil {
		m.setStatus("Error resolving path: "+err.Error(), "red")
		return m, nil
	}

	// Check if file exists, create if it doesn't
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		m.setStatus("Creating new file: "+filepath.Base(absPath), "yellow")
		m.lines = []string{""}
		m.filename = absPath
		m.modified = true
		m.saved = false
	} else {
		// Load the file
		lines, err := loadFile(absPath)
		if err != nil {
			m.setStatus("Error loading file: "+err.Error(), "red")
			return m, nil
		}
		m.lines = lines
		m.filename = absPath
		m.modified = false
		m.saved = true
	}

	// Reset cursor and viewport
	m.cursorX = 0
	m.cursorY = 0
	m.offsetX = 0
	m.offsetY = 0

	// Clear selection
	m.selectionActive = false

	// Invalidate wrap cache
	m.wrapCache = make(map[int][]wrappedLine)

	// Hide file tree and return focus to editor
	m.fileTreeFocused = false
	m.fileTreeVisible = false

	m.setStatus("Opened "+filepath.Base(absPath), "green")
	return m, nil
}

// openFileInNewInstance opens a file in a new TUIWrite instance
// If running in tmux or screen, opens in a split pane
func (m model) openFileInNewInstance(filename string, command string) (tea.Model, tea.Cmd) {
	// Resolve to absolute path
	absPath, err := filepath.Abs(filename)
	if err != nil {
		m.setStatus("Error resolving path: "+err.Error(), "red")
		return m, nil
	}

	// Detect multiplexer environment
	inTmux := os.Getenv("TMUX") != ""
	inScreen := os.Getenv("STY") != ""

	var cmd *exec.Cmd

	if inTmux {
		// Open in tmux split
		splitType := "-h" // horizontal split (side-by-side)
		if command == "split" {
			splitType = "-v" // vertical split (top-bottom)
		}

		cmd = exec.Command("tmux", "split-window", splitType, "-c", filepath.Dir(absPath), "tuiwrite", absPath)
	} else if inScreen {
		// Open in screen split
		cmd = exec.Command("screen", "-X", "split")
		// After split, select the new region and start tuiwrite
		// Note: This is simplified; screen splits are more complex
		if err := cmd.Run(); err == nil {
			cmd = exec.Command("screen", "-X", "focus")
			cmd.Run()
			cmd = exec.Command("screen", "-X", "exec", "tuiwrite", absPath)
		}
	} else {
		// Not in a multiplexer - inform user
		m.setStatus("Not in tmux/screen. Use :e to open in current instance, or start tmux first", "yellow")
		return m, nil
	}

	// Execute the command
	if err := cmd.Start(); err != nil {
		m.setStatus("Failed to open new instance: "+err.Error(), "red")
		return m, nil
	}

	// Detach from the process
	go cmd.Wait()

	m.setStatus("Opened "+filepath.Base(absPath)+" in new instance", "green")
	return m, nil
}

// detectMultiplexer returns the active terminal multiplexer (if any)
func detectMultiplexer() string {
	if os.Getenv("TMUX") != "" {
		return "tmux"
	}
	if os.Getenv("STY") != "" {
		return "screen"
	}
	return ""
}

// getMultiplexerStatus returns a status message about multiplexer availability
func getMultiplexerStatus() string {
	mux := detectMultiplexer()
	switch mux {
	case "tmux":
		return "tmux detected - :new opens in split pane"
	case "screen":
		return "screen detected - :new opens in split pane"
	default:
		return "no multiplexer - use :e to open files"
	}
}

// showMultiplexerHelp displays help about multi-instance commands
func (m model) showMultiplexerHelp() (tea.Model, tea.Cmd) {
	mux := detectMultiplexer()
	var helpMsg string

	if mux == "tmux" {
		helpMsg = "Commands: :e <file> (this window) | :new <file> (side split) | :split <file> (top/bottom)"
	} else if mux == "screen" {
		helpMsg = "Commands: :e <file> (this window) | :new <file> (new split)"
	} else {
		helpMsg = "Commands: :e <file> (open in current instance) | Use tmux/screen for :new"
	}

	m.setStatus(helpMsg, "green")
	return m, nil
}

// getExecutablePath returns the path to the tuiwrite executable
func getExecutablePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(exe)
}

// formatFileSizeForStatus returns a human-readable file size
func formatFileSizeForStatus(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
