package main

import (
	"fmt"
	"os"
	"strings"
)

// TerminalColorSupport represents the level of color support in the terminal
type TerminalColorSupport int

const (
	NoColor TerminalColorSupport = iota
	Color16
	Color256
	TrueColor
)

// checkTerminalColors detects the terminal's color capabilities
func checkTerminalColors() TerminalColorSupport {
	// Check for NO_COLOR environment variable (universal disable)
	if os.Getenv("NO_COLOR") != "" {
		LogInfo("NO_COLOR environment variable detected, disabling colors")
		return NoColor
	}

	// Check COLORTERM for true color support
	colorterm := os.Getenv("COLORTERM")
	if colorterm == "truecolor" || colorterm == "24bit" {
		LogInfo("True color (24-bit) support detected via COLORTERM")
		return TrueColor
	}

	// Check TERM variable
	term := os.Getenv("TERM")

	// True color indicators in TERM
	if strings.Contains(term, "truecolor") || strings.Contains(term, "24bit") {
		LogInfo("True color (24-bit) support detected via TERM")
		return TrueColor
	}

	// 256 color support
	if strings.Contains(term, "256color") {
		LogInfo("256-color support detected, true color recommended for best experience")
		return Color256
	}

	// Basic 16 color support
	if strings.Contains(term, "color") || term == "xterm" || term == "screen" {
		LogWarning("Only 16-color support detected, some colors may not display correctly")
		return Color16
	}

	// Unknown or no color support
	LogWarning("Terminal color support unclear, assuming basic support")
	return Color16
}

// getTerminalInfo returns information about the current terminal
func getTerminalInfo() string {
	term := os.Getenv("TERM")
	colorterm := os.Getenv("COLORTERM")

	info := "Terminal: " + term
	if colorterm != "" {
		info += " (COLORTERM: " + colorterm + ")"
	}

	return info
}

// printColorTest displays all Catppuccin Mocha colors for testing
func printColorTest() {
	println("=== Catppuccin Mocha Color Test ===")
	println()
	println(getTerminalInfo())
	println()

	// Accent colors
	println("ACCENT COLORS:")
	testColor("Rosewater", Rosewater)
	testColor("Flamingo", Flamingo)
	testColor("Pink", Pink)
	testColor("Mauve", Mauve)
	testColor("Red", Red)
	testColor("Maroon", Maroon)
	testColor("Peach", Peach)
	testColor("Yellow", Yellow)
	testColor("Green", Green)
	testColor("Teal", Teal)
	testColor("Sky", Sky)
	testColor("Sapphire", Sapphire)
	testColor("Blue", Blue)
	testColor("Lavender", Lavender)
	println()

	// Monochromatic colors
	println("MONOCHROMATIC COLORS:")
	testColor("Text", Text)
	testColor("Subtext1", Subtext1)
	testColor("Subtext0", Subtext0)
	testColor("Overlay2", Overlay2)
	testColor("Overlay1", Overlay1)
	testColor("Overlay0", Overlay0)
	testColor("Surface2", Surface2)
	testColor("Surface1", Surface1)
	testColor("Surface0", Surface0)
	testColor("Base", Base)
	testColor("Mantle", Mantle)
	testColor("Crust", Crust)
	println()

	println("If colors appear incorrect or monochrome, your terminal may not")
	println("support 24-bit true color. See TERMINAL_SETUP.md for configuration help.")
}

// testColor displays a single color test
func testColor(name string, color Color) {
	// Foreground color
	fg := ColorToANSI(color, true)
	// Background color
	bg := ColorToANSI(color, false)
	// Reset
	reset := "\033[0m"

	// Display: colored text on default bg, and colored background with contrasting text
	fmt.Printf("  %s%-12s%s  %s  %s  %s #%02x%02x%02x\n",
		fg, name, reset,
		bg, "    ", reset,
		color.R, color.G, color.B)
}
