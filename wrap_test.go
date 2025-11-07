package main

import (
	"fmt"
	"testing"
)

// TestLazyWrapping verifies that only visible lines are wrapped
func TestLazyWrapping(t *testing.T) {
	// Create a model with 1000 lines
	m := model{
		lines:     make([]string, 1000),
		width:     80,
		height:    30,
		wrapCache: make(map[int][]wrappedLine),
		cursorX:   0,
		cursorY:   0,
		offsetY:   0,
	}

	// Fill lines with test content
	for i := 0; i < 1000; i++ {
		m.lines[i] = fmt.Sprintf("Line %d: This is a test line with some content that might wrap if it gets long enough.", i)
	}

	// Get visible lines (should only wrap ~28 lines for display)
	visibleHeight := m.height - 2 // 28 lines
	visibleLines := m.getVisibleWrappedLines(0, visibleHeight)

	// Count how many source lines were actually wrapped (cached)
	cachedCount := len(m.wrapCache)

	fmt.Printf("Total lines: %d\n", len(m.lines))
	fmt.Printf("Visible lines requested: %d\n", visibleHeight)
	fmt.Printf("Visible wrapped lines returned: %d\n", len(visibleLines))
	fmt.Printf("Lines actually wrapped (cached): %d\n", cachedCount)

	// With lazy wrapping, we should only have wrapped enough lines to fill the screen
	// This should be much less than 1000
	if cachedCount > 100 {
		t.Errorf("Too many lines wrapped! Expected ~%d, got %d", visibleHeight, cachedCount)
	}

	if cachedCount == 0 {
		t.Error("No lines were wrapped!")
	}

	// Now scroll down and get more visible lines
	m.offsetY = 100
	visibleLines = m.getVisibleWrappedLines(m.offsetY, visibleHeight)
	newCachedCount := len(m.wrapCache)

	fmt.Printf("\nAfter scrolling to offset %d:\n", m.offsetY)
	fmt.Printf("Lines actually wrapped (cached): %d\n", newCachedCount)

	// Cache should have grown
	if newCachedCount <= cachedCount {
		t.Error("Cache didn't grow after scrolling!")
	}

	// The key benefit: calling getVisibleWrappedLines again with the same offset
	// should NOT wrap any more lines (they're all cached now)
	initialCacheSize := len(m.wrapCache)
	visibleLines = m.getVisibleWrappedLines(m.offsetY, visibleHeight)
	finalCacheSize := len(m.wrapCache)

	fmt.Printf("\nAfter second call to getVisibleWrappedLines:\n")
	fmt.Printf("Cache size before: %d, after: %d\n", initialCacheSize, finalCacheSize)

	if finalCacheSize != initialCacheSize {
		t.Errorf("Cache size changed on second call! This means lines were re-wrapped.")
	}

	// Test that editing only invalidates one line
	m.lines[50] = "Modified line!"
	m.invalidateWrapCache(50)

	if len(m.wrapCache) != initialCacheSize-1 {
		t.Errorf("Expected cache size to decrease by 1, got %d", len(m.wrapCache))
	}

	fmt.Printf("\nAfter invalidating line 50:\n")
	fmt.Printf("Cache size: %d (should be %d)\n", len(m.wrapCache), initialCacheSize-1)
}

// TestInvalidateCache verifies cache invalidation works
func TestInvalidateCache(t *testing.T) {
	m := model{
		lines:     []string{"Line 1", "Line 2", "Line 3"},
		width:     80,
		height:    30,
		wrapCache: make(map[int][]wrappedLine),
	}

	// Wrap all lines
	for i := 0; i < len(m.lines); i++ {
		m.getWrappedLine(i)
	}

	if len(m.wrapCache) != 3 {
		t.Errorf("Expected 3 cached lines, got %d", len(m.wrapCache))
	}

	// Invalidate line 1
	m.invalidateWrapCache(1)

	if len(m.wrapCache) != 2 {
		t.Errorf("Expected 2 cached lines after invalidation, got %d", len(m.wrapCache))
	}

	// Verify line 1 is not cached but 0 and 2 are
	if _, exists := m.wrapCache[1]; exists {
		t.Error("Line 1 should not be cached after invalidation")
	}

	if _, exists := m.wrapCache[0]; !exists {
		t.Error("Line 0 should still be cached")
	}

	if _, exists := m.wrapCache[2]; !exists {
		t.Error("Line 2 should still be cached")
	}
}

// TestInvalidateAllCache verifies full cache invalidation
func TestInvalidateAllCache(t *testing.T) {
	m := model{
		lines:     []string{"Line 1", "Line 2", "Line 3"},
		width:     80,
		height:    30,
		wrapCache: make(map[int][]wrappedLine),
	}

	// Wrap all lines
	for i := 0; i < len(m.lines); i++ {
		m.getWrappedLine(i)
	}

	if len(m.wrapCache) != 3 {
		t.Errorf("Expected 3 cached lines, got %d", len(m.wrapCache))
	}

	// Invalidate all
	m.invalidateAllWrapCache()

	if len(m.wrapCache) != 0 {
		t.Errorf("Expected 0 cached lines after full invalidation, got %d", len(m.wrapCache))
	}
}
