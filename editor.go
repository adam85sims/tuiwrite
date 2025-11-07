package main

import tea "github.com/charmbracelet/bubbletea"

// handleEditMode processes keys in edit mode (full editing)
func (m model) handleEditMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up":
		if m.cursorY > 0 {
			m.cursorY--
			// Adjust cursor X if line is shorter
			lineLen := len(m.getCurrentLine())
			if m.cursorX > lineLen {
				m.cursorX = lineLen
			}
			m.adjustViewport()
		}

	case "down":
		if m.cursorY < len(m.lines)-1 {
			m.cursorY++
			// Adjust cursor X if line is shorter
			lineLen := len(m.getCurrentLine())
			if m.cursorX > lineLen {
				m.cursorX = lineLen
			}
			m.adjustViewport()
		}

	case "left":
		if m.cursorX > 0 {
			m.cursorX--
		} else if m.cursorY > 0 {
			// Move to end of previous line
			m.cursorY--
			m.cursorX = len(m.getCurrentLine())
			m.adjustViewport()
		}

	case "right":
		lineLen := len(m.getCurrentLine())
		if m.cursorX < lineLen {
			m.cursorX++
		} else if m.cursorY < len(m.lines)-1 {
			// Move to start of next line
			m.cursorY++
			m.cursorX = 0
			m.adjustViewport()
		}

	case "home":
		m.cursorX = 0

	case "end":
		m.cursorX = len(m.getCurrentLine())

	case "enter":
		// Split line at cursor
		currentLine := m.getCurrentLine()
		before := currentLine[:m.cursorX]
		after := currentLine[m.cursorX:]

		m.lines[m.cursorY] = before
		m.lines = append(m.lines[:m.cursorY+1], append([]string{after}, m.lines[m.cursorY+1:]...)...)

		// Invalidate cache from current line onwards (all subsequent indices shift)
		m.invalidateWrapCacheFrom(m.cursorY)

		m.cursorY++
		m.cursorX = 0
		m.modified = true
		m.adjustViewport()

	case "backspace":
		if m.cursorX > 0 {
			// Delete character before cursor
			line := m.getCurrentLine()
			m.lines[m.cursorY] = line[:m.cursorX-1] + line[m.cursorX:]
			m.cursorX--
			m.modified = true
			m.invalidateWrapCache(m.cursorY) // Invalidate modified line
		} else if m.cursorY > 0 {
			// Join with previous line - this deletes a line, so indices shift
			prevLine := m.lines[m.cursorY-1]
			currentLine := m.getCurrentLine()
			m.lines[m.cursorY-1] = prevLine + currentLine
			m.lines = append(m.lines[:m.cursorY], m.lines[m.cursorY+1:]...)

			// Invalidate cache from previous line onwards (all subsequent indices shift)
			m.invalidateWrapCacheFrom(m.cursorY - 1)

			m.cursorY--
			m.cursorX = len(prevLine)
			m.modified = true
			m.adjustViewport()
		}

	case "delete":
		line := m.getCurrentLine()
		if m.cursorX < len(line) {
			// Delete character at cursor
			m.lines[m.cursorY] = line[:m.cursorX] + line[m.cursorX+1:]
			m.modified = true
			m.invalidateWrapCache(m.cursorY) // Invalidate modified line
		} else if m.cursorY < len(m.lines)-1 {
			// Join with next line - this deletes a line, so indices shift
			nextLine := m.lines[m.cursorY+1]
			m.lines[m.cursorY] = line + nextLine
			m.lines = append(m.lines[:m.cursorY+1], m.lines[m.cursorY+2:]...)

			// Invalidate cache from current line onwards (all subsequent indices shift)
			m.invalidateWrapCacheFrom(m.cursorY)

			m.modified = true
		}

	default:
		// Insert regular characters
		if len(msg.Runes) == 1 {
			r := msg.Runes[0]
			if r >= 32 && r != 127 { // Printable characters
				line := m.getCurrentLine()
				m.lines[m.cursorY] = line[:m.cursorX] + string(r) + line[m.cursorX:]
				m.cursorX++
				m.modified = true
				m.invalidateWrapCache(m.cursorY) // Invalidate modified line
			}
		}
	}

	return m, nil
}
