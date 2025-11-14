#!/bin/bash

# Test script to verify debug logging works

echo "Testing TUIWrite debug logging..."

# Remove old log if exists
rm -f ~/.config/tuiwrite/debug.log

# Create a test file
echo "This is a test file for logging." > /tmp/test-logging.txt

# Run tuiwrite with automated input simulation
# We'll just open and quit immediately to test logging
echo "Opening tuiwrite briefly to generate logs..."
timeout 2s ./tuiwrite /tmp/test-logging.txt || true

# Check if log file was created
if [ -f ~/.config/tuiwrite/debug.log ]; then
    echo ""
    echo "✅ Debug log created successfully!"
    echo ""
    echo "Log location: ~/.config/tuiwrite/debug.log"
    echo ""
    echo "Last 30 lines of debug.log:"
    echo "---"
    tail -30 ~/.config/tuiwrite/debug.log
    echo "---"
    echo ""
    echo "Full log size: $(wc -l ~/.config/tuiwrite/debug.log | awk '{print $1}') lines"
else
    echo "❌ Debug log was not created!"
    exit 1
fi
