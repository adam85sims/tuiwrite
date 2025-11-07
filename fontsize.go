package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Font size constants
const (
	MinFontSize     = 50  // 50% of normal
	MaxFontSize     = 300 // 300% of normal
	DefaultFontSize = 100 // 100% = normal size

	// Size change increments (in percentage points)
	TinyIncrement   = 1 // 6-12pt range
	SmallIncrement  = 2 // 12-18pt range
	MediumIncrement = 4 // 18-24pt range
	LargeIncrement  = 6 // 24+ range

	// Thresholds for different increment sizes (in percentage)
	TinyThreshold   = 80  // Below this: 1% increments
	SmallThreshold  = 120 // Below this: 2% increments
	MediumThreshold = 160 // Below this: 4% increments
	// Above MediumThreshold: 6% increments

	// Timing for held keys
	FontSizeTickInterval = 50 * time.Millisecond // Update every 50ms when held
)

// getFontSizeIncrement returns the appropriate increment based on current size
func getFontSizeIncrement(currentSize int) int {
	switch {
	case currentSize < TinyThreshold:
		return TinyIncrement
	case currentSize < SmallThreshold:
		return SmallIncrement
	case currentSize < MediumThreshold:
		return MediumIncrement
	default:
		return LargeIncrement
	}
}

// increaseFontSize increases the font size by the appropriate increment
func (m *model) increaseFontSize() {
	increment := getFontSizeIncrement(m.fontSize)
	newSize := m.fontSize + increment

	if newSize > MaxFontSize {
		newSize = MaxFontSize
	}

	if newSize != m.fontSize {
		m.fontSize = newSize
		LogDebugf("Font size increased to %d%%", m.fontSize)
		m.setStatus(fmt.Sprintf("Font size: %d%%", m.fontSize), "green")

		// Send ANSI escape sequence to change terminal font size
		// Note: This may not work in all terminals
		m.sendFontSizeEscape()
	}
}

// decreaseFontSize decreases the font size by the appropriate increment
func (m *model) decreaseFontSize() {
	increment := getFontSizeIncrement(m.fontSize)
	newSize := m.fontSize - increment

	if newSize < MinFontSize {
		newSize = MinFontSize
	}

	if newSize != m.fontSize {
		m.fontSize = newSize
		LogDebugf("Font size decreased to %d%%", m.fontSize)
		m.setStatus(fmt.Sprintf("Font size: %d%%", m.fontSize), "green")

		// Send ANSI escape sequence to change terminal font size
		m.sendFontSizeEscape()
	}
}

// sendFontSizeEscape sends terminal escape sequences to change font size
// This uses xterm escape sequences - may not work in all terminals
func (m *model) sendFontSizeEscape() {
	// Calculate font point size from percentage
	// Assuming default is 12pt
	baseSize := 12.0
	pointSize := (float64(m.fontSize) / 100.0) * baseSize

	// Some terminals support these sequences:
	// ESC]50;fontsize=<size><BEL> for xterm
	// For now, we'll just track it internally and let the terminal's
	// native Ctrl+/- handle it if the user's terminal doesn't support our escape

	// We can use this to adjust our rendering if needed
	LogDebugf("Font size set to %.1fpt (%.0f%%)", pointSize, float64(m.fontSize))
}

// tickFontSize returns a command that sends fontSizeTickMsg periodically
func tickFontSize() tea.Cmd {
	return tea.Tick(FontSizeTickInterval, func(t time.Time) tea.Msg {
		return fontSizeTickMsg(t)
	})
}

// handleFontSizeIncrease handles the font size increase (tap or hold)
func (m model) handleFontSizeIncrease() (tea.Model, tea.Cmd) {
	m.increaseFontSize()
	m.fontSizeDirection = "increase"
	m.lastFontSizeTime = time.Now()
	LogEvent("FONT_SIZE", "Increase (hold to continue)")
	return m, tickFontSize()
}

// handleFontSizeDecrease handles the font size decrease (tap or hold)
func (m model) handleFontSizeDecrease() (tea.Model, tea.Cmd) {
	m.decreaseFontSize()
	m.fontSizeDirection = "decrease"
	m.lastFontSizeTime = time.Now()
	LogEvent("FONT_SIZE", "Decrease (hold to continue)")
	return m, tickFontSize()
}

// handleFontSizeTick handles the periodic tick when font size key is held
func (m model) handleFontSizeTick() (tea.Model, tea.Cmd) {
	// Check if enough time has passed and direction is set
	if m.fontSizeDirection != "" {
		elapsed := time.Since(m.lastFontSizeTime)

		// If more than 500ms since last change, assume key was released
		if elapsed > 500*time.Millisecond {
			m.fontSizeDirection = ""
			return m, nil
		}

		// Continue changing size
		if m.fontSizeDirection == "increase" {
			m.increaseFontSize()
		} else if m.fontSizeDirection == "decrease" {
			m.decreaseFontSize()
		}
		m.lastFontSizeTime = time.Now()
		return m, tickFontSize()
	}
	return m, nil
}

// resetFontSize resets to default size
func (m *model) resetFontSize() {
	m.fontSize = DefaultFontSize
	LogEvent("FONT_SIZE", "Reset to default (100%)")
	m.setStatus("Font size reset to 100%", "green")
	m.sendFontSizeEscape()
}
