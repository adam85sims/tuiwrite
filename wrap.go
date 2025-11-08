package main

import "strings"

// getVisibleWrappedLines returns only the wrapped lines needed for the current viewport
// This is the core of the lazy wrapping system - only wraps what's visible!
func (m *model) getVisibleWrappedLines(startOffset int, count int) []wrappedLine {
	if m.width <= 0 {
		return nil
	}

	// Calculate wrap width
	wrapWidth := m.width - 2
	if wrapWidth < 20 {
		wrapWidth = 20
	}

	// If width changed, invalidate all cache
	if m.wrapWidth != wrapWidth {
		m.wrapWidth = wrapWidth
		m.invalidateAllWrapCache()
	}

	// Initialize cache if needed
	if m.wrapCache == nil {
		m.wrapCache = make(map[int][]wrappedLine)
	}

	result := make([]wrappedLine, 0, count)
	wrappedIdx := 0

	// Iterate through source lines and collect wrapped lines
	for lineIdx := 0; lineIdx < len(m.lines) && len(result) < count+startOffset; lineIdx++ {
		// Get wrapped lines for this source line (from cache or wrap it)
		wrappedForLine := m.getWrappedLine(lineIdx)

		// Add wrapped lines to result if they're in the visible range
		for _, wl := range wrappedForLine {
			if wrappedIdx >= startOffset && wrappedIdx < startOffset+count {
				result = append(result, wl)
			}
			wrappedIdx++
			if wrappedIdx >= startOffset+count {
				return result
			}
		}
	}

	return result
}

// getWrappedLine returns wrapped lines for a single source line (cached)
func (m *model) getWrappedLine(lineIdx int) []wrappedLine {
	// Check if already in cache
	if cached, ok := m.wrapCache[lineIdx]; ok {
		return cached
	}

	// Not in cache, wrap it now
	if lineIdx >= len(m.lines) {
		return []wrappedLine{{text: "", sourceLineY: lineIdx, isLastWrap: true}}
	}

	line := m.lines[lineIdx]
	wrappedTexts := wrapLine(line, m.wrapWidth)

	// Convert to wrappedLine structs
	result := make([]wrappedLine, len(wrappedTexts))
	for i, text := range wrappedTexts {
		result[i] = wrappedLine{
			text:        text,
			sourceLineY: lineIdx,
			isLastWrap:  i == len(wrappedTexts)-1,
		}
	}

	// Cache it
	m.wrapCache[lineIdx] = result
	return result
}

// invalidateWrapCache invalidates the cache for a specific line (called when line is edited)
func (m *model) invalidateWrapCache(lineIdx int) {
	if m.wrapCache != nil {
		delete(m.wrapCache, lineIdx)
	}
}

// invalidateWrapCacheFrom invalidates cache for all lines from lineIdx onwards
// This is needed when inserting/deleting lines, as all subsequent line indices shift
func (m *model) invalidateWrapCacheFrom(lineIdx int) {
	if m.wrapCache == nil {
		return
	}

	// Delete all cache entries for lines >= lineIdx
	for i := lineIdx; i < len(m.lines)+10; i++ { // +10 to catch any extras
		delete(m.wrapCache, i)
	}
}

// invalidateAllWrapCache clears the entire wrap cache (called on window resize)
func (m *model) invalidateAllWrapCache() {
	m.wrapCache = make(map[int][]wrappedLine)
}

// rewrapLines is kept for compatibility but now just invalidates cache
// The actual wrapping happens lazily in getVisibleWrappedLines
func (m *model) rewrapLines() {
	if m.width <= 0 {
		return
	}

	// Calculate wrap width
	wrapWidth := m.width - 2
	if wrapWidth < 20 {
		wrapWidth = 20
	}

	// If width changed, invalidate cache
	if m.wrapWidth != wrapWidth {
		m.wrapWidth = wrapWidth
		m.invalidateAllWrapCache()
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
	wrappedIdx := 0

	// Count wrapped lines before cursor's source line
	for lineIdx := 0; lineIdx < m.cursorY && lineIdx < len(m.lines); lineIdx++ {
		wrappedForLine := m.getWrappedLine(lineIdx)
		wrappedIdx += len(wrappedForLine)
	}

	// Now we're at the start of the cursor's source line
	// Find which wrapped line within this source line contains the cursor
	if m.cursorY < len(m.lines) {
		wrappedForLine := m.getWrappedLine(m.cursorY)

		// Determine which wrapped segment contains the cursor based on cursorX
		charCount := 0
		for i, wl := range wrappedForLine {
			lineLen := len(wl.text)
			if m.cursorX <= charCount+lineLen || i == len(wrappedForLine)-1 {
				// Cursor is in this wrapped line
				wrappedIdx += i
				break
			}
			charCount += lineLen
		}
	}

	return wrappedIdx
}

// moveToWrappedLine moves the cursor to a specific wrapped line index
// Returns true if successful, false if out of bounds
func (m *model) moveToWrappedLine(targetWrappedIdx int) bool {
	if targetWrappedIdx < 0 {
		return false
	}

	wrappedIdx := 0

	// Remember current visual X position within the current wrapped line
	currentWrappedIdx := m.getWrappedLineIndexForCursor()
	visualX := m.cursorX

	// Calculate visual X position within current wrapped line
	if m.cursorY < len(m.lines) {
		wrappedForCurrentLine := m.getWrappedLine(m.cursorY)
		charCount := 0
		for _, wl := range wrappedForCurrentLine {
			if currentWrappedIdx == wrappedIdx {
				visualX = m.cursorX - charCount
				break
			}
			wrappedIdx++
			charCount += len(wl.text)
		}
	}

	wrappedIdx = 0

	// Find which source line and wrapped offset contains the target
	for lineIdx := 0; lineIdx < len(m.lines); lineIdx++ {
		wrappedForLine := m.getWrappedLine(lineIdx)

		if wrappedIdx+len(wrappedForLine) > targetWrappedIdx {
			// Target is within this source line
			wrappedOffset := targetWrappedIdx - wrappedIdx

			// Update cursor Y position
			m.cursorY = lineIdx

			// Calculate the starting character position for this wrapped line
			charPos := 0
			for i := 0; i < wrappedOffset && i < len(wrappedForLine); i++ {
				charPos += len(wrappedForLine[i].text)
			}

			// Try to position cursor at similar visual X within this wrapped line
			if wrappedOffset < len(wrappedForLine) {
				wrapLineText := wrappedForLine[wrappedOffset].text
				newCursorX := charPos + visualX

				// Clamp to this wrapped line's actual length
				maxPos := charPos + len(wrapLineText)
				if newCursorX > maxPos {
					newCursorX = maxPos
				}

				// Clamp to source line length
				if newCursorX > len(m.lines[lineIdx]) {
					newCursorX = len(m.lines[lineIdx])
				}

				m.cursorX = newCursorX
			}

			return true
		}

		wrappedIdx += len(wrappedForLine)
	}

	return false
}
