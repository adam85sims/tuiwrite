# TUIWrite API and Extension System

**Last Updated**: November 13, 2025  
**Status**: Planning Document

---

## Overview

This document outlines the API and plugin architecture for TUIWrite, focusing on two primary use cases:
1. **Document Commands** - In-document formatting directives for export processing
2. **Language Extensions** - Syntax highlighting and compilation support for code editing

---

## Use Case 1: Document Commands

### Purpose
Embedded commands within documents (e.g., `#center`, `#pagebreak`, `#font`) that control formatting during export to PDF, HTML, DOCX, etc.

### Design Philosophy
- Commands should be **non-intrusive** in the editor (visible but subtle)
- Commands should be **line-based** for simplicity (e.g., `#center This text is centered`)
- Commands should support **block syntax** for multi-line effects (e.g., `#begin center ... #end center`)
- Export should be **pluggable** - different exporters can interpret commands differently

### Command Syntax Proposals

#### Single-Line Commands
```
#center This text will be centered
#right Align this right
#pagebreak
#image path/to/image.png "Caption text"
```

#### Block Commands
```
#begin center
This entire section
will be centered
#end center

#begin aside
Sidebar content here
#end aside
```

#### Inline Commands
```
This text has #bold(bold text) and #italic(italic text) inline.
Use #color(red)(this text is red).
```

### Command Categories

**Formatting**:
- `#center`, `#right`, `#left`, `#justify` - Text alignment
- `#bold`, `#italic`, `#underline` - Text styling (inline)
- `#font(name)`, `#size(12pt)` - Font control
- `#color(name)` - Text color

**Structure**:
- `#pagebreak` - Force page break in PDF/print exports
- `#section`, `#chapter` - Semantic document structure
- `#begin`/`#end` - Block containers

**Media**:
- `#image(path)` - Insert image
- `#video(path)` - Embed video (for HTML export)
- `#audio(path)` - Embed audio

**Metadata**:
- `#title`, `#author`, `#date` - Document metadata
- `#tag(keyword)` - Tagging for organization

**Screenplay-Specific** (when in script mode):
- `#scene`, `#action`, `#dialogue`, `#character` - Fountain-style markup
- `#transition`, `#parenthetical` - Standard screenplay elements

### Implementation Strategy

#### Phase 1: Parser
```go
// Command represents a parsed document command
type Command struct {
    Type      string            // e.g., "center", "pagebreak", "image"
    Args      []string          // positional arguments
    Options   map[string]string // named options
    Content   string            // for inline commands
    LineStart int               // source line where command starts
    LineEnd   int               // source line where command ends (for blocks)
    IsBlock   bool              // true for #begin/#end pairs
}

// CommandParser extracts commands from document
type CommandParser struct {
    lines []string
}

func (cp *CommandParser) ParseCommands() []Command
func (cp *CommandParser) IsCommandLine(line string) bool
func (cp *CommandParser) ParseCommand(line string) (Command, error)
```

#### Phase 2: Visual Rendering
- Commands should be visually distinct in the editor
- Suggestion: Dim gray text with special icon/color
- Don't break reading flow
- Add to `view.go` rendering pipeline

```go
// In view.go
func (m model) renderLineWithCommands(line string) string {
    if m.isCommandLine(line) {
        // Render with special styling
        commandStyle := lipgloss.NewStyle().
            Foreground(lipgloss.Color(ColorToHex(Overlay0))).
            Italic(true)
        return commandStyle.Render(line)
    }
    return line
}
```

#### Phase 3: Export System
```go
// Exporter interface for different output formats
type Exporter interface {
    Export(content []string, commands []Command) ([]byte, error)
    SupportsCommand(cmdType string) bool
}

// PDF exporter using commands
type PDFExporter struct {
    commands []Command
}

func (pe *PDFExporter) Export(content []string, commands []Command) ([]byte, error)

// HTML exporter
type HTMLExporter struct {
    commands []Command
}

func (he *HTMLExporter) Export(content []string, commands []Command) ([]byte, error)

// Markdown exporter (might strip commands or convert to markdown)
type MarkdownExporter struct{}
```

#### Phase 4: Export Commands
```
:export pdf output.pdf
:export html output.html
:export docx output.docx
:export md output.md
```

---

## Use Case 2: Language Extensions

### Purpose
Support syntax highlighting and optional compilation/execution for code editing use cases.

### Design Philosophy
- Extensions should be **language-agnostic** - define a standard interface
- Syntax highlighting should be **fast and lazy** - only highlight visible lines
- Compilation/execution should be **optional** - not everyone needs it
- Extensions should be **sandboxed** - no arbitrary code execution from extensions

### Architecture Proposal

#### Extension Definition Format (TOML)

Extensions live in `~/.config/tuiwrite/extensions/` as self-contained directories:

```
~/.config/tuiwrite/extensions/
├── python/
│   ├── extension.toml
│   ├── syntax.json
│   └── compile.sh (optional)
├── go/
│   ├── extension.toml
│   ├── syntax.json
│   └── compile.sh (optional)
└── rust/
    ├── extension.toml
    ├── syntax.json
    └── compile.sh (optional)
```

**extension.toml**:
```toml
[extension]
name = "Python"
version = "1.0.0"
author = "TUIWrite"
description = "Python syntax highlighting and execution"

[language]
file_extensions = [".py", ".pyw"]
comment_line = "#"
comment_block_start = "'''"
comment_block_end = "'''"

[features]
syntax_highlighting = true
compilation = false
execution = true

[execution]
command = "python3"
args = ["{file}"]
```

**syntax.json** (TextMate-style grammar subset):
```json
{
  "name": "Python",
  "scopeName": "source.python",
  "patterns": [
    {
      "name": "keyword.control.python",
      "match": "\\b(if|elif|else|for|while|break|continue|return|def|class|import|from|as|try|except|finally|with|raise)\\b"
    },
    {
      "name": "string.quoted.double.python",
      "begin": "\"",
      "end": "\""
    },
    {
      "name": "comment.line.python",
      "match": "#.*$"
    },
    {
      "name": "constant.numeric.python",
      "match": "\\b\\d+(\\.\\d+)?\\b"
    }
  ],
  "colors": {
    "keyword.control": "Mauve",
    "string.quoted": "Green",
    "comment.line": "Overlay0",
    "constant.numeric": "Peach"
  }
}
```

### Implementation Strategy

#### Phase 1: Extension Manager
```go
// Extension represents a loaded language extension
type Extension struct {
    Name        string
    Version     string
    Language    LanguageConfig
    Features    FeatureSet
    Syntax      *SyntaxHighlighter
    Execution   *ExecutionConfig
}

type LanguageConfig struct {
    FileExtensions     []string
    CommentLine        string
    CommentBlockStart  string
    CommentBlockEnd    string
}

type FeatureSet struct {
    SyntaxHighlighting bool
    Compilation        bool
    Execution          bool
}

type ExtensionManager struct {
    extensions map[string]*Extension // keyed by language name
    extensionsDir string
}

func NewExtensionManager() *ExtensionManager
func (em *ExtensionManager) LoadExtensions() error
func (em *ExtensionManager) GetExtensionForFile(filename string) *Extension
func (em *ExtensionManager) GetExtensionByName(name string) *Extension
```

#### Phase 2: Syntax Highlighting Engine
```go
// SyntaxHighlighter applies syntax highlighting to lines
type SyntaxHighlighter struct {
    patterns []SyntaxPattern
    colors   map[string]Color
}

type SyntaxPattern struct {
    Name  string
    Regex *regexp.Regexp
    Color Color
}

func (sh *SyntaxHighlighter) HighlightLine(line string) string {
    // Apply patterns and return line with ANSI/lipgloss styling
    // Similar to how spell-check highlighting works
}

// Integration in view.go
func (m model) renderLineWithSyntax(line string) string {
    if m.currentExtension != nil && m.currentExtension.Features.SyntaxHighlighting {
        return m.currentExtension.Syntax.HighlightLine(line)
    }
    return line
}
```

#### Phase 3: Compilation/Execution
```go
type ExecutionConfig struct {
    Command string
    Args    []string
}

func (ec *ExecutionConfig) Execute(filename string) (string, error) {
    // Substitute {file} in args with filename
    // Run command with exec
    // Return output
}

// New command mode commands:
// :run - Execute current file with extension's execution config
// :compile - Compile current file (if supported)
```

#### Phase 4: Extension Commands
```
:extension list                  # List installed extensions
:extension load python           # Explicitly load Python extension
:extension reload                # Reload all extensions
:syntax on|off                   # Toggle syntax highlighting
:run                            # Run current file with extension's executor
```

---

## Security Considerations

### Document Commands
- **Low Risk**: Commands are declarative, not executable
- **Validation**: Sanitize paths in `#image`, `#video`, etc.
- **Exporter Safety**: Exporters should escape user content appropriately

### Language Extensions
- **Sandbox Execution**: Run compilation/execution in isolated processes
- **Resource Limits**: Set timeouts and memory limits for execution
- **No Arbitrary Code**: Extensions use declarative config (TOML/JSON), not executable scripts
  - Exception: Optional `compile.sh` scripts must be explicitly trusted by user
- **User Confirmation**: Prompt before first execution of any file
- **Extension Verification**: Consider signed extensions or curated repository

---

## API Surface

### Core Interfaces

```go
// DocumentProcessor handles command parsing and export
type DocumentProcessor interface {
    ParseCommands(lines []string) []Command
    Export(format string, output string) error
}

// LanguageSupport handles syntax and execution
type LanguageSupport interface {
    GetName() string
    HighlightLine(line string) string
    CanExecute() bool
    Execute(filename string) (string, error)
    CanCompile() bool
    Compile(filename string) error
}

// Plugin interface for future extensibility
type Plugin interface {
    Name() string
    Version() string
    Initialize(editor *Editor) error
    OnDocumentLoad(filename string)
    OnDocumentSave(filename string)
    OnCommand(cmd string, args []string) error
}
```

### Configuration

**~/.config/tuiwrite/config.toml**:
```toml
[editor]
enable_commands = true
enable_extensions = true

[commands]
render_style = "subtle"  # "subtle", "hidden", "normal"
export_default_format = "pdf"

[extensions]
enabled = ["python", "go", "rust"]
auto_detect = true  # Auto-load extension based on file extension

[security]
allow_execution = true
execution_timeout = 30  # seconds
prompt_before_execution = true
```

---

## Migration Path

### Short Term (v0.2)
1. Implement basic command parser for `#pagebreak`, `#center`, etc.
2. Add visual distinction for command lines in editor
3. Implement simple PDF exporter using command hints
4. Add `:export` command

### Medium Term (v0.3)
1. Create extension manager framework
2. Implement syntax highlighting engine
3. Ship with 3-5 built-in extensions (Python, Go, JavaScript, Markdown, JSON)
4. Add `:run` and `:syntax` commands

### Long Term (v1.0)
1. Public extension API documentation
2. Extension marketplace/repository
3. Advanced features: LSP integration, autocomplete, jump-to-definition
4. Plugin system for arbitrary editor enhancements

---

## Examples

### Example 1: Story with Formatting Commands

```markdown
#title My Short Story
#author Adam Sims
#date 2025-11-13

#begin center
**Chapter 1**
The Beginning
#end center

Once upon a time, in a land far away...

#pagebreak

#begin center
**Chapter 2**
The Middle
#end center

The adventure continued...
```

Export to PDF would use these commands to format the document appropriately.

### Example 2: Code Editing with Python Extension

```python
# File: script.py (automatically loads Python extension)

def hello_world():
    print("Hello from TUIWrite!")
    
if __name__ == "__main__":
    hello_world()
```

Commands available:
- `:run` - Executes `python3 script.py`
- `:syntax on` - Enables Python syntax highlighting
- Keywords, strings, and numbers are color-coded automatically

---

## Open Questions

1. **Command Conflict Resolution**: What if a user wants `#` to be literal in their document?
   - Solution: Escape with `\#` or add setting to disable commands
   
2. **Extension Distribution**: How do users discover and install extensions?
   - Solution: Built-in extension manager with GitHub repository backend
   
3. **Performance**: Will syntax highlighting on every keystroke be too slow?
   - Solution: Use lazy highlighting (only visible lines) and debouncing
   
4. **Extension Updates**: How to handle breaking changes in extension format?
   - Solution: Version the extension API, maintain backward compatibility
   
5. **Multi-Language Files**: What about files with embedded languages (HTML with JS/CSS)?
   - Solution: Phase 2 feature - nested syntax highlighting

---

## References

- TextMate Grammar: https://macromates.com/manual/en/language_grammars
- Tree-sitter: https://tree-sitter.github.io/tree-sitter/
- LSP Protocol: https://microsoft.github.io/language-server-protocol/
- Fountain Spec: https://fountain.io/syntax

---

**End of API Document**
