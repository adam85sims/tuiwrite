# Files Need Restoration

**Date**: November 7, 2025  
**Issue**: wrap.go and view.go are using old `wrappedLines` field instead of new `wrapCache` implementation

## What Happened

During the tab refactoring attempt, we ran `git restore` on several files. This worked for most files, but wrap.go and view.go had never been committed with the lazy wrapping implementation, so they reverted to the old inefficient version.

## What Needs To Be Fixed

### wrap.go
- Currently uses: `m.wrappedLines` (old array of all wrapped lines)
- Should use: `m.wrapCache` (map of cached wrapped lines per source line)
- Key functions needed:
  - `getVisibleWrappedLines(offsetY, count)` - Returns only visible wrapped lines
  - `wrapSingleLine(lineIdx)` - Wraps one line and caches it
  - `invalidateWrapCache(lineIdx)` - Clears cache for edited line
  - `getTotalWrappedLineCount()` - Counts total wrapped lines
  - `getWrappedLineIndexForCursor()` - Maps cursor position to wrapped line index

### view.go  
- Currently uses: `m.wrappedLines[i]` (accessing pre-wrapped array)
- Should use: `m.getVisibleWrappedLines(m.offsetY, visibleHeight)` (on-demand wrapping)

## Reference Files

The correct implementation is documented in:
- `LAZY_WRAPPING.md` - Full explanation and benchmark results
- `wrap_test.go` - Test cases showing correct usage
- `benchmark_test.go` - Performance benchmarks

The implementation achieved **3,192x speedup** by wrapping only visible lines.

## Quick Fix Options

### Option 1: Recreate from documentation
Use LAZY_WRAPPING.md as a guide to rewrite wrap.go and view.go

### Option 2: Ask AI to regenerate
Provide LAZY_WRAPPING.md context and ask to implement the functions listed above

### Option 3: Manual recovery
If you have backups from earlier today (before git restore), those would have the correct versions

## Current Build Status

```
ERROR: m.wrappedLines undefined (should be m.wrapCache)
```

Files affected:
- wrap.go (7 references to wrappedLines)
- view.go (5 references to wrappedLines)

Once these are fixed to use wrapCache and the lazy wrapping functions, the build will succeed.
