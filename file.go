package main

import (
	"os"
	"strings"
	"time"
)

// saveFile writes the document to disk
func (m *model) saveFile() error {
	content := strings.Join(m.lines, "\n")
	err := os.WriteFile(m.filename, []byte(content), 0644)
	if err != nil {
		return err
	}

	m.modified = false
	m.lastSave = time.Now()
	return nil
}

// loadFile reads the document from disk
func loadFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// Create new empty file
			return []string{""}, nil
		}
		return nil, err
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	// Ensure at least one line exists
	if len(lines) == 0 {
		lines = []string{""}
	}

	return lines, nil
}
