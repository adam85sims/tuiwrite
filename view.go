package main

import (
	"fmt"
	"strings"
	"time"
)

// View renders the UI
func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	var sb strings.Builder

	// Calculate visible area (leave 2 lines for status bar)
	visibleHeight := m.height - 2

	// Render document lines using wrapped lines
	for i := 0; i < visibleHeight; i++ {
		wrappedIdx := m.offsetY + i

		if wrappedIdx < len(m.wrappedLines) {
			wl := m.wrappedLines[wrappedIdx]
			line := wl.text

			// Show cursor if this wrapped line is from the current source line
			// Cursor is visible in both Read and Edit modes
			if wl.sourceLineY == m.cursorY {
				// Calculate which part of the source line is in this wrap
				wrapStartPos := 0
				for j := 0; j < wrappedIdx; j++ {
					if m.wrappedLines[j].sourceLineY == m.cursorY {
						wrapStartPos += len(m.wrappedLines[j].text)
						if j+1 < len(m.wrappedLines) && m.wrappedLines[j+1].sourceLineY == m.cursorY {
							wrapStartPos++ // Account for space removed during wrapping
						}
					}
				}

				// Check if cursor is in this wrapped line segment
				wrapEndPos := wrapStartPos + len(line)
				if m.cursorX >= wrapStartPos && m.cursorX <= wrapEndPos {
					cursorPosInWrap := m.cursorX - wrapStartPos

					if cursorPosInWrap <= len(line) {
						before := line[:cursorPosInWrap]
						after := ""
						cursor := "█"
						if cursorPosInWrap < len(line) {
							cursor = string(line[cursorPosInWrap])
							after = line[cursorPosInWrap+1:]
						}

						// Different cursor style for Read vs Edit mode
						if m.mode == ReadMode {
							// Underlined cursor for read mode
							sb.WriteString(before + "\033[4m" + cursor + "\033[0m" + after)
						} else {
							// Inverted cursor for edit mode
							sb.WriteString(before + "\033[7m" + cursor + "\033[0m" + after)
						}
					} else {
						if m.mode == ReadMode {
							sb.WriteString(line + "\033[4m \033[0m")
						} else {
							sb.WriteString(line + "\033[7m \033[0m")
						}
					}
				} else {
					sb.WriteString(line)
				}
			} else {
				sb.WriteString(line)
			}
		} else {
			sb.WriteString("~") // Empty line indicator
		}

		if i < visibleHeight-1 {
			sb.WriteString("\n")
		}
	}

	// Status bar (2 lines)
	sb.WriteString("\n")
	sb.WriteString(m.renderStatusBar())

	return sb.String()
}

// renderStatusBar creates the status bar display
func (m model) renderStatusBar() string {
	modifiedIndicator := ""
	if m.modified {
		modifiedIndicator = " [+]"
	}

	// First line: mode, filename, position
	leftStatus := fmt.Sprintf(" %s | %s%s", m.mode, m.filename, modifiedIndicator)
	rightStatus := fmt.Sprintf("Ln %d, Col %d ", m.cursorY+1, m.cursorX+1)

	padding := m.width - len(leftStatus) - len(rightStatus)
	if padding < 0 {
		padding = 0
	}

	statusLine1 := "\033[7m" + leftStatus + strings.Repeat(" ", padding) + rightStatus + "\033[0m"

	// Second line: status message or help
	statusLine2 := ""
	if m.statusMsg.Text != "" && time.Since(m.statusMsg.Timestamp) < 3*time.Second {
		// Show temporary status message
		colorCode := ""
		switch m.statusMsg.Color {
		case "green":
			colorCode = "\033[32m"
		case "yellow":
			colorCode = "\033[33m"
		case "red":
			colorCode = "\033[31m"
		}
		statusLine2 = colorCode + m.statusMsg.Text + "\033[0m"
	} else if m.commandMode {
		// Show command buffer when in command mode
		statusLine2 = m.commandBuffer
	} else {
		// Show help text
		if m.mode == ReadMode {
			statusLine2 = "Press INSERT to edit | : for commands | Ctrl+S to save | Ctrl+C to quit"
		} else {
			statusLine2 = "Press ESC for read mode | Ctrl+S to save | Ctrl+C to quit"
		}
	}

	return statusLine1 + "\n" + statusLine2
}

// setStatus sets a status message with color
func (m *model) setStatus(text string, color string) {
	m.statusMsg = StatusMessage{
		Text:      text,
		Color:     color,
		Timestamp: time.Now(),
	}
}
