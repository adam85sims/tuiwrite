package main

import tea "github.com/charmbracelet/bubbletea"

// handleReadMode processes keys in read mode (navigation only)
func (m model) handleReadMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursorY > 0 {
			m.cursorY--
			m.adjustViewport()
		}

	case "down", "j":
		if m.cursorY < len(m.lines)-1 {
			m.cursorY++
			m.adjustViewport()
		}

	case "left", "h":
		if m.cursorX > 0 {
			m.cursorX--
		}

	case "right", "l":
		if m.cursorX < len(m.getCurrentLine()) {
			m.cursorX++
		}

	case "pgup":
		m.cursorY -= (m.height - 4) // Leave room for status bar
		if m.cursorY < 0 {
			m.cursorY = 0
		}
		m.adjustViewport()

	case "pgdown":
		m.cursorY += (m.height - 4)
		if m.cursorY >= len(m.lines) {
			m.cursorY = len(m.lines) - 1
		}
		m.adjustViewport()

	case "ctrl+home", "g":
		if msg.String() == "g" {
			// Simple implementation: single 'g' goes to top
			m.cursorY = 0
			m.cursorX = 0
			m.adjustViewport()
		} else {
			m.cursorY = 0
			m.cursorX = 0
			m.adjustViewport()
		}

	case "ctrl+end", "G":
		m.cursorY = len(m.lines) - 1
		m.cursorX = len(m.getCurrentLine())
		m.adjustViewport()

	case "home":
		m.cursorX = 0

	case "end":
		m.cursorX = len(m.getCurrentLine())
	}

	return m, nil
}

// getCurrentLine returns the current line content
func (m model) getCurrentLine() string {
	if m.cursorY >= 0 && m.cursorY < len(m.lines) {
		return m.lines[m.cursorY]
	}
	return ""
}

// adjustViewport ensures cursor is visible
func (m *model) adjustViewport() {
	visibleHeight := m.height - 2

	// Find the wrapped line index where cursor is located
	cursorWrappedIdx := m.getWrappedLineIndexForCursor()

	// Scroll down if cursor is below viewport
	if cursorWrappedIdx >= m.offsetY+visibleHeight {
		m.offsetY = cursorWrappedIdx - visibleHeight + 1
	}

	// Scroll up if cursor is above viewport
	if cursorWrappedIdx < m.offsetY {
		m.offsetY = cursorWrappedIdx
	}
}
