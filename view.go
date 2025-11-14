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

			// Show cursor if this wrapped line is from the current source line
			if wl.sourceLineY == m.cursorY && !m.fileTreeFocused {
				// Apply both spell-check and cursor
				line = m.renderLineWithCursorAndSpellCheck(line, m.cursorX)
			} else {
				// Just apply spell-check highlighting (no cursor)
				line = m.applySpellCheckHighlighting(line)
			}

			// Apply selection highlighting if active
			if m.selectionActive {
				line = m.applySelectionHighlighting(line, wl.sourceLineY, 0)
			}

			// Truncate line if too long for editor width (after styling)
				// Note: This is a rough truncation and may cut ANSI codes, but it's acceptable
				if len(line) > editorWidth-1 {
					line = line[:editorWidth-2] + "…"
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
				// Apply both spell-check and cursor highlighting
				line = m.renderLineWithCursorAndSpellCheck(line, m.cursorX)
			} else {
				// Just apply spell-check highlighting (no cursor)
				line = m.applySpellCheckHighlighting(line)
			}

			// Apply selection highlighting if active
			if m.selectionActive {
				line = m.applySelectionHighlighting(line, wl.sourceLineY, 0)
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
			commandText = "Press INSERT to edit | : for commands | Ctrl+S to save | Ctrl+Q to quit"
		} else {
			commandText = "ESC: read | Shift+arrows: select | Ctrl+C: copy | Ctrl+V: paste | Ctrl+S: save | Ctrl+Q: quit"
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

// applySpellCheckHighlighting applies red background to misspelled words in a line
func (m model) applySpellCheckHighlighting(line string) string {
	// Skip if spell checker is disabled
	if m.spellChecker == nil || !m.spellChecker.enabled {
		return line
	}

	// Get words and their positions
	words := getWordsInLine(line)
	if len(words) == 0 {
		return line
	}

	// Build the highlighted line by processing each segment
	var result strings.Builder
	lastEnd := 0

	// Style for misspelled words - red background with base (dark) text
	errorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(ColorToHex(Red))).
		Foreground(lipgloss.Color(ColorToHex(Base)))

	for _, w := range words {
		// Add text before this word
		if w.start > lastEnd {
			result.WriteString(line[lastEnd:w.start])
		}

		// Check if word is misspelled
		if !m.spellChecker.checkWord(w.word) {
			// Highlight misspelled word
			result.WriteString(errorStyle.Render(w.word))
		} else {
			// Word is correct, render normally
			result.WriteString(w.word)
		}

		lastEnd = w.end
	}

	// Add any remaining text after the last word
	if lastEnd < len(line) {
		result.WriteString(line[lastEnd:])
	}

	return result.String()
}

// renderLineWithCursorAndSpellCheck renders a line with both cursor and spell-check highlighting
func (m model) renderLineWithCursorAndSpellCheck(line string, cursorPos int) string {
	// First, work out which character will have the cursor
	if cursorPos < 0 {
		cursorPos = 0
	}
	
	// For cursor positioning, we need to work with the original line
	// We'll build the line with spell-check and insert cursor styling at the right position
	
	// Get words for spell-check
	words := getWordsInLine(line)
	
	// Styles
	errorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(ColorToHex(Red))).
		Foreground(lipgloss.Color(ColorToHex(Base)))
	
	cursorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(ColorToHex(Maroon))).
		Foreground(lipgloss.Color(ColorToHex(Base)))
	
	// Build result
	var result strings.Builder
	lastEnd := 0
	spellCheckEnabled := m.spellChecker != nil && m.spellChecker.enabled
	
	for _, w := range words {
		// Render text before this word
		if w.start > lastEnd {
			segment := line[lastEnd:w.start]
			// Check if cursor is in this segment
			if cursorPos >= lastEnd && cursorPos < w.start {
				relPos := cursorPos - lastEnd
				result.WriteString(segment[:relPos])
				if m.cursorVisible {
					cursorChar := " "
					if relPos < len(segment) {
						cursorChar = string(segment[relPos])
					}
					result.WriteString(cursorStyle.Render(cursorChar))
				} else {
					if relPos < len(segment) {
						result.WriteString(string(segment[relPos]))
					} else {
						result.WriteString(" ")
					}
				}
				if relPos+1 < len(segment) {
					result.WriteString(segment[relPos+1:])
				}
			} else {
				result.WriteString(segment)
			}
		}
		
		// Render the word
		isMisspelled := spellCheckEnabled && !m.spellChecker.checkWord(w.word)
		
		// Check if cursor is in this word
		if cursorPos >= w.start && cursorPos < w.end {
			relPos := cursorPos - w.start
			before := w.word[:relPos]
			cursorChar := string(w.word[relPos])
			after := ""
			if relPos+1 < len(w.word) {
				after = w.word[relPos+1:]
			}
			
			if isMisspelled {
				result.WriteString(errorStyle.Render(before))
			} else {
				result.WriteString(before)
			}
			
			if m.cursorVisible {
				result.WriteString(cursorStyle.Render(cursorChar))
			} else {
				result.WriteString(cursorChar)
			}
			
			if isMisspelled {
				result.WriteString(errorStyle.Render(after))
			} else {
				result.WriteString(after)
			}
		} else {
			// No cursor in word
			if isMisspelled {
				result.WriteString(errorStyle.Render(w.word))
			} else {
				result.WriteString(w.word)
			}
		}
		
		lastEnd = w.end
	}
	
	// Render remaining text after last word
	if lastEnd < len(line) {
		segment := line[lastEnd:]
		if cursorPos >= lastEnd && cursorPos < len(line) {
			relPos := cursorPos - lastEnd
			result.WriteString(segment[:relPos])
			if m.cursorVisible {
				result.WriteString(cursorStyle.Render(string(segment[relPos])))
			} else {
				result.WriteString(string(segment[relPos]))
			}
			if relPos+1 < len(segment) {
				result.WriteString(segment[relPos+1:])
			}
		} else {
			result.WriteString(segment)
		}
	}
	
	// Cursor at end of line
	if cursorPos >= len(line) && m.cursorVisible {
		result.WriteString(cursorStyle.Render(" "))
	}
	
	return result.String()
}

// applySelectionHighlighting applies blue highlight background to selected text
func (m model) applySelectionHighlighting(line string, lineY int, lineOffset int) string {
	if !m.selectionActive {
		return line
	}

	// Get selection bounds (normalize so start is always before end)
	startY, startX := m.selectionStartY, m.selectionStartX
	endY, endX := m.cursorY, m.cursorX

	// Swap if selection is backwards
	if startY > endY || (startY == endY && startX > endX) {
		startY, endY = endY, startY
		startX, endX = endX, startX
	}

	// Check if this line is in the selection range
	if lineY < startY || lineY > endY {
		return line // Not in selection
	}

	// Determine selection range for this line
	var selStart, selEnd int
	if lineY == startY && lineY == endY {
		// Single line selection
		selStart = startX - lineOffset
		selEnd = endX - lineOffset
	} else if lineY == startY {
		// First line of multi-line selection
		selStart = startX - lineOffset
		selEnd = len(line)
	} else if lineY == endY {
		// Last line of multi-line selection
		selStart = 0
		selEnd = endX - lineOffset
	} else {
		// Middle line of multi-line selection
		selStart = 0
		selEnd = len(line)
	}

	// Clamp to line bounds
	if selStart < 0 {
		selStart = 0
	}
	if selEnd > len(line) {
		selEnd = len(line)
	}
	if selStart >= len(line) {
		return line // Selection doesn't overlap this wrapped segment
	}

	// This is a simplified approach - it won't handle ANSI codes properly
	// A full implementation would need to parse and re-apply ANSI codes
	// For now, we'll apply selection over the plain text
	selectionStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(ColorToHex(Blue))).
		Foreground(lipgloss.Color(ColorToHex(Base)))

	if selStart == 0 && selEnd >= len(line) {
		// Entire line is selected
		return selectionStyle.Render(line)
	} else if selStart > 0 && selEnd >= len(line) {
		// Selection starts partway through
		return line[:selStart] + selectionStyle.Render(line[selStart:])
	} else if selStart == 0 {
		// Selection ends partway through
		return selectionStyle.Render(line[:selEnd]) + line[selEnd:]
	} else {
		// Selection in middle of line
		return line[:selStart] + selectionStyle.Render(line[selStart:selEnd]) + line[selEnd:]
	}
}
