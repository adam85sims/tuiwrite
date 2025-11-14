package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// buildFileTree scans the directory and builds the file tree structure
func buildFileTree(rootPath string) ([]FileNode, error) {
	// Get absolute path
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	// Read root directory
	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, err
	}

	// Separate folders and files
	var folders []FileNode
	var files []FileNode

	for _, entry := range entries {
		// Skip hidden files and common non-document directories
		name := entry.Name()
		if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" {
			continue
		}

		fullPath := filepath.Join(absPath, name)

		node := FileNode{
			Name:     name,
			Path:     fullPath,
			IsDir:    entry.IsDir(),
			Expanded: false,
			Depth:    0,
		}

		// For directories, recursively get children (but don't expand by default)
		if entry.IsDir() {
			children, _ := getDirectoryChildren(fullPath, 1)
			node.Children = children
			folders = append(folders, node)
		} else {
			// Only include text-like files
			if isTextFile(name) {
				files = append(files, node)
			}
		}
	}

	// Sort folders and files alphabetically
	sort.Slice(folders, func(i, j int) bool {
		return strings.ToLower(folders[i].Name) < strings.ToLower(folders[j].Name)
	})
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	// Combine: folders first, then files
	result := append(folders, files...)

	LogInfof("Built file tree with %d folders and %d files", len(folders), len(files))
	return result, nil
}

// getDirectoryChildren recursively gets children for a directory
func getDirectoryChildren(dirPath string, depth int) ([]FileNode, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var folders []FileNode
	var files []FileNode

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" {
			continue
		}

		fullPath := filepath.Join(dirPath, name)

		node := FileNode{
			Name:     name,
			Path:     fullPath,
			IsDir:    entry.IsDir(),
			Expanded: false,
			Depth:    depth,
		}

		if entry.IsDir() {
			// Get children for subdirectories
			children, _ := getDirectoryChildren(fullPath, depth+1)
			node.Children = children
			folders = append(folders, node)
		} else {
			if isTextFile(name) {
				files = append(files, node)
			}
		}
	}

	// Sort folders and files
	sort.Slice(folders, func(i, j int) bool {
		return strings.ToLower(folders[i].Name) < strings.ToLower(folders[j].Name)
	})
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	return append(folders, files...), nil
}

// isTextFile checks if a file is likely a text file we want to show
func isTextFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	textExtensions := []string{
		".txt", ".md", ".markdown",
		".rst", ".org", ".adoc",
		".tex", ".fountain",
		".go", ".py", ".js", ".ts",
		".html", ".css", ".json", ".yaml", ".yml",
		".c", ".cpp", ".h", ".java",
		".sh", ".bash", ".zsh",
	}

	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}

	// Also accept files without extension
	return ext == ""
}

// flattenFileTree creates a flat list of visible nodes for rendering
func flattenFileTree(nodes []FileNode) []FileNode {
	var result []FileNode

	for _, node := range nodes {
		result = append(result, node)

		// If directory is expanded, add its children
		if node.IsDir && node.Expanded {
			result = append(result, flattenFileTree(node.Children)...)
		}
	}

	return result
}

// toggleFileTree shows/hides the file tree sidebar
func (m model) toggleFileTree() (tea.Model, tea.Cmd) {
	m.fileTreeVisible = !m.fileTreeVisible

	if m.fileTreeVisible {
		// Transfer focus to file tree when opening
		m.fileTreeFocused = true
		LogEvent("FILE_TREE", "Opened file tree (focus transferred)")
		m.setStatus("File tree - Use arrows to navigate, Enter to select, F1 to close", "green")
	} else {
		// Return focus to editor when closing
		m.fileTreeFocused = false
		LogEvent("FILE_TREE", "Closed file tree (focus returned to editor)")
		m.setStatus("File tree closed", "yellow")
	}

	return m, nil
}

// handleFileTreeNavigation handles navigation within the file tree
func (m model) handleFileTreeNavigation(key string) (tea.Model, tea.Cmd) {
	if !m.fileTreeVisible || !m.fileTreeFocused {
		return m, nil
	}

	flatNodes := flattenFileTree(m.fileTreeNodes)
	visibleHeight := m.height - 2 // Leave room for status bar

	switch key {
	case "up", "k":
		if m.fileTreeCursor > 0 {
			m.fileTreeCursor--
			// Scroll up if cursor goes above viewport
			if m.fileTreeCursor < m.fileTreeOffset {
				m.fileTreeOffset = m.fileTreeCursor
			}
		}

	case "down", "j":
		if m.fileTreeCursor < len(flatNodes)-1 {
			m.fileTreeCursor++
			// Scroll down if cursor goes below viewport
			if m.fileTreeCursor >= m.fileTreeOffset+visibleHeight {
				m.fileTreeOffset = m.fileTreeCursor - visibleHeight + 1
			}
		}

	case "enter":
		// Toggle directory or open file
		if m.fileTreeCursor < len(flatNodes) {
			node := flatNodes[m.fileTreeCursor]

			if node.IsDir {
				// Toggle directory expansion
				m.toggleDirectory(node.Path)
				LogDebugf("Toggled directory: %s", node.Name)
			} else {
				// Open file for editing in current instance
				LogInfof("Opening file: %s", node.Path)
				return m.openFileInCurrentInstance(node.Path)
			}
		}
	}

	return m, nil
}

// toggleDirectory expands or collapses a directory in the tree
func (m *model) toggleDirectory(path string) {
	m.toggleDirectoryRecursive(&m.fileTreeNodes, path)
}

// toggleDirectoryRecursive recursively finds and toggles a directory
func (m *model) toggleDirectoryRecursive(nodes *[]FileNode, path string) bool {
	for i := range *nodes {
		node := &(*nodes)[i]

		if node.Path == path && node.IsDir {
			node.Expanded = !node.Expanded
			LogDebugf("Directory %s now expanded=%v", node.Name, node.Expanded)
			return true
		}

		// Search in children
		if node.IsDir && len(node.Children) > 0 {
			if m.toggleDirectoryRecursive(&node.Children, path) {
				return true
			}
		}
	}
	return false
}

// initFileTree initializes the file tree for the current working directory
func (m *model) initFileTree() error {
	// Get the directory of the current file
	dir := filepath.Dir(m.filename)
	if dir == "" || dir == "." {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	m.fileTreeRoot = dir

	// Build the file tree
	nodes, err := buildFileTree(dir)
	if err != nil {
		LogErrorf("Failed to build file tree: %v", err)
		return err
	}

	m.fileTreeNodes = nodes
	m.fileTreeCursor = 0

	LogInfof("Initialized file tree for directory: %s", dir)
	return nil
}
