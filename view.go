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

	// If file tree is visible, render split view
	if m.fileTreeVisible {
		// File tree takes 30 characters, editor gets the rest
		treeWidth := 30
		editorWidth := m.width - treeWidth - 1 // -1 for divider

		// Get flattened file tree nodes
		flatNodes := flattenFileTree(m.fileTreeNodes)

		// Get visible lines once for the editor (not in the loop!)
		visibleLines := m.getVisibleWrappedLines(m.offsetY, visibleHeight)

		// Render each line with file tree on left, editor on right
		for i := 0; i < visibleHeight; i++ {
			// Calculate the actual node index with offset
			nodeIdx := m.fileTreeOffset + i

			// Render file tree column
			if nodeIdx < len(flatNodes) {
				node := flatNodes[nodeIdx]

				// Highlight selected node
				prefix := " "
				suffix := ""
				if nodeIdx == m.fileTreeCursor && m.fileTreeFocused {
					prefix = "\033[7m" // Inverted
					suffix = "\033[0m"
				}

				// Indentation
				indent := strings.Repeat("  ", node.Depth)

				// Icon
				icon := "📄 "
				if node.IsDir {
					if node.Expanded {
						icon = "📂 "
					} else {
						icon = "📁 "
					}
				}

				// Truncate name if too long
				name := node.Name
				maxNameLen := treeWidth - len(indent) - len(icon) - 2
				if len(name) > maxNameLen {
					name = name[:maxNameLen-1] + "…"
				}

				treeLine := prefix + indent + icon + name + suffix

				// Pad to tree width
				padding := treeWidth - len(stripAnsi(treeLine))
				if padding > 0 {
					treeLine += strings.Repeat(" ", padding)
				}

				sb.WriteString(treeLine)
			} else {
				// Empty tree line
				sb.WriteString(strings.Repeat(" ", treeWidth))
			}

			// Divider
			sb.WriteString("│")

			// Render editor column
			if i < len(visibleLines) {
				wl := visibleLines[i]
				line := wl.text

				// Truncate line if too long for editor width
				if len(line) > editorWidth-1 {
					line = line[:editorWidth-2] + "…"
				}

				// Show cursor if this wrapped line is from the current source line
				if wl.sourceLineY == m.cursorY && !m.fileTreeFocused {
					cursorPosInWrap := m.cursorX
					if cursorPosInWrap > len(line) {
						cursorPosInWrap = len(line)
					}

					before := ""
					after := ""
					cursor := "█"

					if cursorPosInWrap < len(line) {
						before = line[:cursorPosInWrap]
						cursor = string(line[cursorPosInWrap])
						after = line[cursorPosInWrap+1:]
					} else if cursorPosInWrap == len(line) {
						before = line
					}

					if m.mode == ReadMode {
						sb.WriteString(before + "\033[4m" + cursor + "\033[0m" + after)
					} else {
						sb.WriteString(before + "\033[7m" + cursor + "\033[0m" + after)
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
	} else {
		// No file tree - full width editor
		// Get only the wrapped lines we need for the visible area (lazy!)
		visibleLines := m.getVisibleWrappedLines(m.offsetY, visibleHeight)

		// Render document lines using wrapped lines
		for i := 0; i < visibleHeight; i++ {
			if i < len(visibleLines) {
				wl := visibleLines[i]
				line := wl.text

				// Show cursor if this wrapped line is from the current source line
				// Cursor is visible in both Read and Edit modes
				if wl.sourceLineY == m.cursorY {
					// Calculate which part of the source line is in this wrap
					// For now, simplified - just show cursor at the position
					// TODO: Handle cursor position within wrapped lines more accurately
					cursorPosInWrap := m.cursorX
					if cursorPosInWrap > len(line) {
						cursorPosInWrap = len(line)
					}

					before := ""
					after := ""
					cursor := "█"

					if cursorPosInWrap < len(line) {
						before = line[:cursorPosInWrap]
						cursor = string(line[cursorPosInWrap])
						after = line[cursorPosInWrap+1:]
					} else if cursorPosInWrap == len(line) {
						before = line
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
					sb.WriteString(line)
				}
			} else {
				sb.WriteString("~") // Empty line indicator
			}

			if i < visibleHeight-1 {
				sb.WriteString("\n")
			}
		}
	}

	// Status bar (2 lines)
	sb.WriteString("\n")
	sb.WriteString(m.renderStatusBar())

	return sb.String()
}

// stripAnsi removes ANSI escape sequences for length calculation
func stripAnsi(s string) string {
	// Simple implementation - removes common ANSI codes
	s = strings.ReplaceAll(s, "\033[7m", "")
	s = strings.ReplaceAll(s, "\033[0m", "")
	s = strings.ReplaceAll(s, "\033[4m", "")
	return s
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
