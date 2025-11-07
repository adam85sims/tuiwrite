package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Logger handles debug logging to file
type Logger struct {
	file    *os.File
	enabled bool
}

var globalLogger *Logger

// initLogger initializes the debug logger
func initLogger() error {
	logger := &Logger{
		enabled: true, // Can be toggled via command-line flag later
	}

	if !logger.enabled {
		globalLogger = logger
		return nil
	}

	// Get config directory (same as dictionaries)
	configDir, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Open or create debug.log file
	logPath := filepath.Join(configDir, "debug.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	logger.file = file
	globalLogger = logger

	// Log session start
	logger.log("SESSION_START", "TUIWrite started")
	logger.logf("INFO", "Log file: %s", logPath)
	logger.logf("INFO", "Platform: %s/%s", runtime.GOOS, runtime.GOARCH)

	return nil
}

// getConfigDir returns the platform-specific config directory
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "tuiwrite"), nil
		}
		return filepath.Join(homeDir, "AppData", "Roaming", "tuiwrite"), nil
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", "tuiwrite"), nil
	default: // Linux and other Unix-like systems
		return filepath.Join(homeDir, ".config", "tuiwrite"), nil
	}
}

// log writes a log entry with timestamp
func (l *Logger) log(level string, message string) {
	if !l.enabled || l.file == nil {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	logLine := fmt.Sprintf("[%s] %s: %s\n", timestamp, level, message)
	l.file.WriteString(logLine)
}

// logf writes a formatted log entry
func (l *Logger) logf(level string, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.log(level, message)
}

// Close closes the log file
func (l *Logger) Close() error {
	if l.file != nil {
		l.log("SESSION_END", "TUIWrite exiting")
		return l.file.Close()
	}
	return nil
}

// Helper functions for global logger access

// LogDebug logs a debug message
func LogDebug(message string) {
	if globalLogger != nil {
		globalLogger.log("DEBUG", message)
	}
}

// LogDebugf logs a formatted debug message
func LogDebugf(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.logf("DEBUG", format, args...)
	}
}

// LogInfo logs an info message
func LogInfo(message string) {
	if globalLogger != nil {
		globalLogger.log("INFO", message)
	}
}

// LogInfof logs a formatted info message
func LogInfof(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.logf("INFO", format, args...)
	}
}

// LogWarning logs a warning message
func LogWarning(message string) {
	if globalLogger != nil {
		globalLogger.log("WARNING", message)
	}
}

// LogWarningf logs a formatted warning message
func LogWarningf(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.logf("WARNING", format, args...)
	}
}

// LogError logs an error message
func LogError(message string) {
	if globalLogger != nil {
		globalLogger.log("ERROR", message)
	}
}

// LogErrorf logs a formatted error message
func LogErrorf(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.logf("ERROR", format, args...)
	}
}

// LogEvent logs a specific event (mode changes, commands, etc.)
func LogEvent(event string, details string) {
	if globalLogger != nil {
		globalLogger.logf("EVENT", "%s: %s", event, details)
	}
}

// CloseLogger closes the global logger
func CloseLogger() {
	if globalLogger != nil {
		globalLogger.Close()
	}
}
