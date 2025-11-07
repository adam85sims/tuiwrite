package main

import tea "github.com/charmbracelet/bubbletea"

// toggleZenMode toggles between normal and zen (fullscreen/borderless) mode
func (m model) toggleZenMode() (tea.Model, tea.Cmd) {
	m.zenMode = !m.zenMode

	if m.zenMode {
		LogEvent("ZEN_MODE", "Entered zen mode (fullscreen)")
		m.setStatus("Zen Mode: ON (F11 to exit)", "green")
		// Enter fullscreen/borderless mode
		return m, tea.EnterAltScreen
	} else {
		LogEvent("ZEN_MODE", "Exited zen mode")
		m.setStatus("Zen Mode: OFF", "green")
		// We're already in alt screen, but we can send a message
		return m, nil
	}
}

// getZenModeIndicator returns a visual indicator for zen mode in the status bar
func (m model) getZenModeIndicator() string {
	if m.zenMode {
		return " [ZEN]"
	}
	return ""
}
