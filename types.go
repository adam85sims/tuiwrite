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

	// Word wrap (lazy caching)
	wrapCache map[int][]wrappedLine // cache of wrapped lines per source line index
	wrapWidth int                   // width used for current wrapping

	// Font size controls
	fontSize          int       // current font size (100 = 100% = normal)
	fontSizeDirection string    // "increase", "decrease", or "" when not held
	lastFontSizeTime  time.Time // for smooth font size changes

	// Zen mode
	zenMode bool // true when in fullscreen/borderless mode

	// Spell checking
	spellChecker *SpellChecker

	// Command mode
	commandMode   bool   // true when in command mode (after typing :)
	commandBuffer string // current command being typed

	// File tree
	fileTreeVisible bool       // true when file tree sidebar is shown
	fileTreeFocused bool       // true when file tree has focus (vs editor)
	fileTreeCursor  int        // current selection in file tree
	fileTreeOffset  int        // scroll offset for file tree viewport
	fileTreeNodes   []FileNode // flattened list of visible tree nodes
	fileTreeRoot    string     // root directory path for file tree
}

// FileNode represents an item in the file tree
type FileNode struct {
	Name     string     // file or folder name
	Path     string     // full path to the file/folder
	IsDir    bool       // true if this is a directory
	Expanded bool       // true if directory is expanded (only applies to directories)
	Depth    int        // indentation depth in the tree
	Children []FileNode // child nodes (only for directories)
}

// Tab represents a single open file/document (for future multi-tab support)
type Tab struct {
	Filename string   // file path
	DocMode  DocMode  // story or script mode
	Lines    []string // document content
	CursorX  int      // column position
	CursorY  int      // line position
	OffsetY  int      // vertical scroll offset
	OffsetX  int      // horizontal scroll offset
	Modified bool     // has unsaved changes
	LastSave time.Time
}

// autoSaveMsg is sent periodically to trigger auto-save
type autoSaveMsg time.Time

// fontSizeTickMsg is sent periodically when font size key is held
type fontSizeTickMsg time.Time
