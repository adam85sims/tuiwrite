package main

import (
	"strings"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

// handleEditMode processes keys in edit mode (full editing)
func (m model) handleEditMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "shift+up", "shift+down", "shift+left", "shift+right":
		// Start selection if not already active
		if !m.selectionActive {
			m.selectionActive = true
			m.selectionStartX = m.cursorX
			m.selectionStartY = m.cursorY
		}

		// Move cursor based on direction (selection remains anchored at start)
		switch msg.String() {
		case "shift+up":
			currentWrappedIdx := m.getWrappedLineIndexForCursor()
			if currentWrappedIdx > 0 {
				m.moveToWrappedLine(currentWrappedIdx - 1)
				lineLen := len(m.getCurrentLine())
				if m.cursorX > lineLen {
					m.cursorX = lineLen
				}
				m.adjustViewport()
			}
		case "shift+down":
			currentWrappedIdx := m.getWrappedLineIndexForCursor()
			m.moveToWrappedLine(currentWrappedIdx + 1)
			lineLen := len(m.getCurrentLine())
			if m.cursorX > lineLen {
				m.cursorX = lineLen
			}
			m.adjustViewport()
		case "shift+left":
			if m.cursorX > 0 {
				m.cursorX--
			} else if m.cursorY > 0 {
				m.cursorY--
				m.cursorX = len(m.getCurrentLine())
				m.adjustViewport()
			}
		case "shift+right":
			lineLen := len(m.getCurrentLine())
			if m.cursorX < lineLen {
				m.cursorX++
			} else if m.cursorY < len(m.lines)-1 {
				m.cursorY++
				m.cursorX = 0
				m.adjustViewport()
			}
		}

	case "up":
		// Clear selection on non-shift movement
		m.selectionActive = false
		// Navigate by wrapped lines, not source lines
		currentWrappedIdx := m.getWrappedLineIndexForCursor()
		if currentWrappedIdx > 0 {
			m.moveToWrappedLine(currentWrappedIdx - 1)
			// Adjust cursor X if line is shorter
			lineLen := len(m.getCurrentLine())
			if m.cursorX > lineLen {
				m.cursorX = lineLen
			}
			m.adjustViewport()
		}

	case "down":
		// Clear selection on non-shift movement
		m.selectionActive = false
		// Navigate by wrapped lines, not source lines
		currentWrappedIdx := m.getWrappedLineIndexForCursor()
		m.moveToWrappedLine(currentWrappedIdx + 1)
		// Adjust cursor X if line is shorter
		lineLen := len(m.getCurrentLine())
		if m.cursorX > lineLen {
			m.cursorX = lineLen
		}
		m.adjustViewport()

	case "left":
		// Clear selection on non-shift movement
		m.selectionActive = false
		if m.cursorX > 0 {
			m.cursorX--
		} else if m.cursorY > 0 {
			// Move to end of previous line
			m.cursorY--
			m.cursorX = len(m.getCurrentLine())
			m.adjustViewport()
		}

	case "right":
		// Clear selection on non-shift movement
		m.selectionActive = false
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
		// Clear selection on non-shift movement
		m.selectionActive = false
		m.cursorX = 0

	case "end":
		// Clear selection on non-shift movement
		m.selectionActive = false
		m.cursorX = len(m.getCurrentLine())

	case "tab":
		// Insert 4 spaces for tab
		line := m.getCurrentLine()
		m.lines[m.cursorY] = line[:m.cursorX] + "    " + line[m.cursorX:]
		m.cursorX += 4
		m.modified = true
		m.invalidateWrapCache(m.cursorY)

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

	case "ctrl+c":
		// Copy selected text to clipboard
		if !m.selectionActive {
			m.setStatus("No text selected", "yellow")
			return m, nil
		}

		// Get selection bounds (normalize so start is always before end)
		startY, startX := m.selectionStartY, m.selectionStartX
		endY, endX := m.cursorY, m.cursorX

		// Swap if selection is backwards
		if startY > endY || (startY == endY && startX > endX) {
			startY, endY = endY, startY
			startX, endX = endX, startX
		}

		var selectedText string
		if startY == endY {
			// Single line selection
			selectedText = m.lines[startY][startX:endX]
		} else {
			// Multi-line selection
			var textBuilder strings.Builder
			// First line: from startX to end
			textBuilder.WriteString(m.lines[startY][startX:])
			textBuilder.WriteString("\n")

			// Middle lines: entire lines
			for i := startY + 1; i < endY; i++ {
				textBuilder.WriteString(m.lines[i])
				textBuilder.WriteString("\n")
			}

			// Last line: from start to endX
			textBuilder.WriteString(m.lines[endY][:endX])
			selectedText = textBuilder.String()
		}

		// Copy to clipboard
		err := clipboard.WriteAll(selectedText)
		if err != nil {
			m.setStatus("Failed to copy to clipboard", "red")
		} else {
			m.setStatus("Copied to clipboard", "green")
		}

		// Keep selection active after copy (user might want to see what was copied)
		return m, nil

	case "ctrl+v":
		// Paste from clipboard
		clipText, err := clipboard.ReadAll()
		if err != nil {
			m.setStatus("Failed to read clipboard", "red")
			return m, nil
		}

		if clipText == "" {
			return m, nil // Nothing to paste
		}

		// Split pasted text into lines
		pasteLines := strings.Split(clipText, "\n")

		if len(pasteLines) == 1 {
			// Single line paste - insert at cursor position
			line := m.getCurrentLine()
			m.lines[m.cursorY] = line[:m.cursorX] + clipText + line[m.cursorX:]
			m.cursorX += len(clipText)
			m.modified = true
			m.invalidateWrapCache(m.cursorY)
		} else {
			// Multi-line paste
			currentLine := m.getCurrentLine()
			before := currentLine[:m.cursorX]
			after := currentLine[m.cursorX:]

			// First line: before + first paste line
			m.lines[m.cursorY] = before + pasteLines[0]

			// Insert middle lines
			newLines := make([]string, 0, len(pasteLines)-1)
			for i := 1; i < len(pasteLines)-1; i++ {
				newLines = append(newLines, pasteLines[i])
			}

			// Last line: last paste line + after
			lastPasteLine := pasteLines[len(pasteLines)-1]
			newLines = append(newLines, lastPasteLine+after)

			// Insert all new lines after current line
			m.lines = append(m.lines[:m.cursorY+1], append(newLines, m.lines[m.cursorY+1:]...)...)

			// Move cursor to end of pasted text
			m.cursorY += len(pasteLines) - 1
			m.cursorX = len(lastPasteLine)

			// Invalidate cache from original line onwards (lines were inserted)
			m.invalidateWrapCacheFrom(m.cursorY - len(pasteLines) + 1)

			m.modified = true
			m.adjustViewport()
		}

	default:
		// Insert regular characters
		if len(msg.Runes) == 1 {
			r := msg.Runes[0]
			if r >= 32 && r != 127 { // Printable characters
				// Clear selection on typing
				m.selectionActive = false
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
