package main

import "time"

// Mode represents the current editing mode
type Mode int

const (
	ReadMode Mode = iota
	EditMode
)

func (m Mode) String() string {
	switch m {
	case ReadMode:
		return "READ"
	case EditMode:
		return "EDIT"
	default:
		return "UNKNOWN"
	}
}

// DocMode represents the document type
type DocMode string

const (
	StoryMode  DocMode = "story"
	ScriptMode DocMode = "script"
)

// StatusMessage represents a temporary message to display
type StatusMessage struct {
	Text      string
	Color     string // "green", "yellow", "red"
	Timestamp time.Time
}

// wrappedLine represents a line after word wrapping
type wrappedLine struct {
	text        string // the wrapped line text
	sourceLineY int    // which original line this came from
	isLastWrap  bool   // true if this is the last wrap of the source line
}

// model represents the application state
type model struct {
	// File information
	filename string
	docMode  DocMode

	// Document content
	lines []string

	// Cursor position
	cursorX int // column position (in actual line)
	cursorY int // line position (actual line number)

	// View state
	offsetY int // vertical scroll offset (in wrapped lines)
	offsetX int // horizontal scroll offset

	// Mode state
	mode Mode

	// Status and messages
	statusMsg StatusMessage
	saved     bool
	modified  bool

	// Dimensions
	width  int
	height int

	// Auto-save
	lastSave time.Time

	// Word wrap
	wrappedLines []wrappedLine // cached wrapped lines for display
	wrapWidth    int           // width used for current wrapping

	// Spell checking
	spellChecker *SpellChecker

	// Command mode
	commandMode   bool   // true when in command mode (after typing :)
	commandBuffer string // current command being typed
}

// autoSaveMsg is sent periodically to trigger auto-save
type autoSaveMsg time.Time
