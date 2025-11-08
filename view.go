package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// View renders the UI
func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	var sb strings.Builder

	// Base style with Catppuccin Base background and Text foreground
	baseStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(ColorToHex(Base))).
		Foreground(lipgloss.Color(ColorToHex(Text)))

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

			var treeLine string
			// Render file tree column
			if nodeIdx < len(flatNodes) {
				node := flatNodes[nodeIdx]

				// Color based on file type
				var itemColor Color
				if node.IsDir {
					itemColor = Yellow // Folders are yellow
				} else {
					itemColor = Blue // Files are blue
				}

				// Highlight selected node with inverted colors
				var treeStyle lipgloss.Style
				if nodeIdx == m.fileTreeCursor && m.fileTreeFocused {
					// Inverted: yellow/blue background with base text
					treeStyle = lipgloss.NewStyle().
						Background(lipgloss.Color(ColorToHex(itemColor))).
						Foreground(lipgloss.Color(ColorToHex(Base)))
				} else {
					// Normal: base background with colored text
					treeStyle = lipgloss.NewStyle().
						Background(lipgloss.Color(ColorToHex(Base))).
						Foreground(lipgloss.Color(ColorToHex(itemColor)))
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

				treeLine = treeStyle.Width(treeWidth).Render(indent + icon + name)
			} else {
				// Empty tree line with base background
				treeLine = baseStyle.Width(treeWidth).Render("")
			}

			sb.WriteString(treeLine)

			// Divider with base background
			divider := baseStyle.Render("│")
			sb.WriteString(divider)

			// Render editor column
			var editorLine string
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
					// Bounds checking - prevent negative indices and out of range
					if cursorPosInWrap < 0 {
						cursorPosInWrap = 0
					}
					if cursorPosInWrap > len(line) {
						cursorPosInWrap = len(line)
					}

					before := ""
					after := ""
					cursor := " " // Default to space for empty position

					if cursorPosInWrap < len(line) {
						before = line[:cursorPosInWrap]
						cursor = string(line[cursorPosInWrap])
						after = line[cursorPosInWrap+1:]
					} else if cursorPosInWrap == len(line) {
						before = line
					}

					// Apply Maroon background to cursor if visible
					if m.cursorVisible {
						cursorStyle := lipgloss.NewStyle().
							Background(lipgloss.Color(ColorToHex(Maroon))).
							Foreground(lipgloss.Color(ColorToHex(Base)))
						line = before + cursorStyle.Render(cursor) + after
					} else {
						// Cursor not visible - just show the character normally
						line = before + cursor + after
					}
				}

				editorLine = baseStyle.Width(editorWidth).Render(line)
			} else {
				editorLine = baseStyle.Width(editorWidth).Render("~") // Empty line indicator
			}

			sb.WriteString(editorLine)

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
			var line string
			if i < len(visibleLines) {
				wl := visibleLines[i]
				line = wl.text

				// Show cursor if this wrapped line is from the current source line
				// Cursor is visible in both Read and Edit modes
				if wl.sourceLineY == m.cursorY {
					// Calculate which part of the source line is in this wrap
					// For now, simplified - just show cursor at the position
					// TODO: Handle cursor position within wrapped lines more accurately
					cursorPosInWrap := m.cursorX
					// Bounds checking - prevent negative indices and out of range
					if cursorPosInWrap < 0 {
						cursorPosInWrap = 0
					}
					if cursorPosInWrap > len(line) {
						cursorPosInWrap = len(line)
					}

					before := ""
					after := ""
					cursor := " " // Default to space for empty position

					if cursorPosInWrap < len(line) {
						before = line[:cursorPosInWrap]
						cursor = string(line[cursorPosInWrap])
						after = line[cursorPosInWrap+1:]
					} else if cursorPosInWrap == len(line) {
						before = line
					}

					// Apply Maroon background to cursor if visible
					if m.cursorVisible {
						cursorStyle := lipgloss.NewStyle().
							Background(lipgloss.Color(ColorToHex(Maroon))).
							Foreground(lipgloss.Color(ColorToHex(Base)))
						line = before + cursorStyle.Render(cursor) + after
					} else {
						// Cursor not visible - just show the character normally
						line = before + cursor + after
					}
				}
			} else {
				line = "~" // Empty line indicator
			}

			// Apply base style with full width to ensure background fills
			styledLine := baseStyle.Width(m.width).Render(line)
			sb.WriteString(styledLine)

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
	// Status bar styles
	statusStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(ColorToHex(Surface0))).
		Foreground(lipgloss.Color(ColorToHex(Text))).
		Width(m.width)

	commandStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(ColorToHex(Mantle))).
		Foreground(lipgloss.Color(ColorToHex(Text))).
		Width(m.width)

	modifiedIndicator := ""
	if m.modified {
		modifiedIndicator = " [+]"
	}

	// First line: mode, filename, position (Surface0 background)
	leftStatus := fmt.Sprintf(" %s | %s%s", m.mode, m.filename, modifiedIndicator)
	rightStatus := fmt.Sprintf("Ln %d, Col %d ", m.cursorY+1, m.cursorX+1)

	padding := m.width - len(leftStatus) - len(rightStatus)
	if padding < 0 {
		padding = 0
	}

	statusLine1 := statusStyle.Render(leftStatus + strings.Repeat(" ", padding) + rightStatus)

	// Second line: status message or help (Mantle background for visual separation)
	var commandText string
	if m.statusMsg.Text != "" && time.Since(m.statusMsg.Timestamp) < 3*time.Second {
		// Show temporary status message with appropriate color
		var msgColor Color
		switch m.statusMsg.Color {
		case "green":
			msgColor = Green
		case "yellow":
			msgColor = Yellow
		case "red":
			msgColor = Red
		default:
			msgColor = Text
		}
		// Use lipgloss for consistent color rendering
		msgStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorToHex(msgColor)))
		commandText = msgStyle.Render(m.statusMsg.Text)
	} else if m.commandMode {
		// Show command buffer when in command mode
		commandText = m.commandBuffer
	} else {
		// Show help text
		if m.mode == ReadMode {
			commandText = "Press INSERT to edit | : for commands | Ctrl+S to save | Ctrl+C to quit"
		} else {
			commandText = "Press ESC for read mode | Ctrl+S to save | Ctrl+C to quit"
		}
	}

	statusLine2 := commandStyle.Render(commandText)

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
