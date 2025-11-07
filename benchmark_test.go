package main

import (
	"testing"
)

// BenchmarkOldApproach simulates the old approach: wrap all lines on every edit
func BenchmarkOldApproachWrappingAllLines(b *testing.B) {
	m := model{
		lines:     make([]string, 10000),
		width:     80,
		height:    30,
		wrapCache: make(map[int][]wrappedLine),
	}

	// Fill with test content
	for i := 0; i < 10000; i++ {
		m.lines[i] = "This is a test line with some content that might wrap."
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate old approach: wrap everything
		m.invalidateAllWrapCache()
		for j := 0; j < len(m.lines); j++ {
			m.getWrappedLine(j)
		}
	}
}

// BenchmarkLazyWrappingVisible simulates the new approach: only wrap visible lines
func BenchmarkLazyWrappingVisibleOnly(b *testing.B) {
	m := model{
		lines:     make([]string, 10000),
		width:     80,
		height:    30,
		wrapCache: make(map[int][]wrappedLine),
	}

	// Fill with test content
	for i := 0; i < 10000; i++ {
		m.lines[i] = "This is a test line with some content that might wrap."
	}

	visibleHeight := m.height - 2

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// New approach: only get visible lines
		m.getVisibleWrappedLines(0, visibleHeight)
	}
}

// BenchmarkEditWithCacheInvalidation simulates editing a single line
func BenchmarkEditWithCacheInvalidation(b *testing.B) {
	m := model{
		lines:     make([]string, 10000),
		width:     80,
		height:    30,
		wrapCache: make(map[int][]wrappedLine),
	}

	// Fill with test content
	for i := 0; i < 10000; i++ {
		m.lines[i] = "This is a test line with some content that might wrap."
	}

	// Pre-populate cache
	visibleHeight := m.height - 2
	m.getVisibleWrappedLines(0, visibleHeight)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lineToEdit := i % 10000
		// Simulate edit: invalidate one line and re-wrap it
		m.invalidateWrapCache(lineToEdit)
		m.lines[lineToEdit] = m.lines[lineToEdit] + "x"
		m.getWrappedLine(lineToEdit)
	}
}
