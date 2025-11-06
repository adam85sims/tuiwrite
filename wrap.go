package main

import "strings"

// rewrapLines recalculates word wrapping for all lines
func (m *model) rewrapLines() {
	if m.width <= 0 {
		return
	}

	// Calculate usable width (leave some margin for status info)
	wrapWidth := m.width - 2
	if wrapWidth < 20 {
		wrapWidth = 20 // Minimum wrap width
	}

	m.wrapWidth = wrapWidth
	m.wrappedLines = make([]wrappedLine, 0)

	for lineIdx, line := range m.lines {
		wrapped := wrapLine(line, wrapWidth)
		for i, wText := range wrapped {
			m.wrappedLines = append(m.wrappedLines, wrappedLine{
				text:        wText,
				sourceLineY: lineIdx,
				isLastWrap:  i == len(wrapped)-1,
			})
		}
	}

	// Adjust viewport to keep cursor visible
	m.adjustViewport()
}

// wrapLine wraps a single line to the specified width
func wrapLine(line string, width int) []string {
	if len(line) == 0 {
		return []string{""}
	}

	if len(line) <= width {
		return []string{line}
	}

	var wrapped []string
	words := strings.Fields(line)

	if len(words) == 0 {
		// Line is only whitespace
		return []string{line}
	}

	currentLine := ""

	for _, word := range words {
		// If word itself is longer than width, break it
		if len(word) > width {
			if currentLine != "" {
				wrapped = append(wrapped, currentLine)
				currentLine = ""
			}
			// Break long word across lines
			for len(word) > width {
				wrapped = append(wrapped, word[:width])
				word = word[width:]
			}
			if len(word) > 0 {
				currentLine = word
			}
			continue
		}

		// Try adding word to current line
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine) <= width {
			currentLine = testLine
		} else {
			// Word doesn't fit, start new line
			if currentLine != "" {
				wrapped = append(wrapped, currentLine)
			}
			currentLine = word
		}
	}

	// Add remaining text
	if currentLine != "" {
		wrapped = append(wrapped, currentLine)
	}

	if len(wrapped) == 0 {
		return []string{""}
	}

	return wrapped
}

// getWrappedLineIndexForCursor returns the wrapped line index for the cursor position
func (m *model) getWrappedLineIndexForCursor() int {
	for i, wl := range m.wrappedLines {
		if wl.sourceLineY == m.cursorY {
			return i
		}
	}
	return 0
}
