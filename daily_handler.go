package logpy

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// DailyFileHandler is a handler that rotates log files daily
type DailyFileHandler struct {
	*baseHandler
	baseDir       string
	filePrefix    string
	dateLayout    string
	maxDaysToKeep int
	currentDate   string
	currentFile   *os.File
	fileMutex     sync.Mutex
	useColor      bool
	colorConfig   ColorConfig
}

// NewDailyFileHandler creates a new daily rotating file handler
// baseDir: directory where log files will be stored (e.g., "./logs")
// filePrefix: optional prefix for log files (e.g., "app" -> "app-2025-11-06.log", empty -> "2025-11-06.log")
// level: minimum log level to handle
// maxDaysToKeep: number of days to retain old log files (0 = keep all)
// useColor: whether to include color codes in the log files
// colorConfig: color configuration for different log levels
func NewDailyFileHandler(baseDir, filePrefix string, level Level, maxDaysToKeep int, useColor bool, colorConfig ColorConfig) (*DailyFileHandler, error) {
	// Use default date layout (ISO 8601)
	dateLayout := "2006-01-02"

	// Create base directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create formatter based on color preference
	var formatter Formatter
	if useColor {
		formatter = &ConsoleFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			AddCaller:       true,
			UseColor:        true,
			ColorConfig:     colorConfig,
		}
	} else {
		formatter = &ConsoleFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			AddCaller:       true,
			UseColor:        false,
			ColorConfig:     colorConfig,
		}
	}

	h := &DailyFileHandler{
		baseDir:       baseDir,
		filePrefix:    filePrefix,
		dateLayout:    dateLayout,
		maxDaysToKeep: maxDaysToKeep,
		useColor:      useColor,
		colorConfig:   colorConfig,
		baseHandler: &baseHandler{
			level:     level,
			formatter: formatter,
		},
	}

	// Open the initial file for today
	if err := h.rotateIfNeeded(); err != nil {
		return nil, err
	}

	// Set the writer to self (we implement io.Writer)
	h.baseHandler.writer = h

	return h, nil
}

// Write implements io.Writer interface with daily rotation
func (h *DailyFileHandler) Write(p []byte) (n int, err error) {
	h.fileMutex.Lock()
	defer h.fileMutex.Unlock()

	// Check if we need to rotate to a new day's file
	if err := h.rotateIfNeeded(); err != nil {
		return 0, err
	}

	// Write to the current file
	return h.currentFile.Write(p)
}

// rotateIfNeeded checks if the date has changed and opens a new file if needed
func (h *DailyFileHandler) rotateIfNeeded() error {
	today := time.Now().Format(h.dateLayout)

	// If we're already on the correct date and file is open, no rotation needed
	if h.currentDate == today && h.currentFile != nil {
		return nil
	}

	// Close the current file if it exists
	if h.currentFile != nil {
		if err := h.currentFile.Close(); err != nil {
			// Log the error but continue with rotation
			fmt.Fprintf(os.Stderr, "error closing log file: %v\n", err)
		}
	}

	// Build the new filename
	filename := h.buildFilename(today)

	// Create the file (append mode, create if doesn't exist)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file %s: %w", filename, err)
	}

	h.currentFile = f
	h.currentDate = today

	// Cleanup old files if configured
	if h.maxDaysToKeep > 0 {
		// Run cleanup in background to avoid blocking
		go h.cleanupOldFiles()
	}

	return nil
}

// buildFilename constructs the full path to the log file for a given date
func (h *DailyFileHandler) buildFilename(date string) string {
	var filename string
	if h.filePrefix != "" {
		filename = h.filePrefix + "-" + date + ".log"
	} else {
		filename = date + ".log"
	}
	return filepath.Join(h.baseDir, filename)
}

// cleanupOldFiles removes log files older than maxDaysToKeep days
func (h *DailyFileHandler) cleanupOldFiles() {
	cutoffDate := time.Now().AddDate(0, 0, -h.maxDaysToKeep)

	files, err := os.ReadDir(h.baseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading log directory for cleanup: %v\n", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Only process .log files
		if filepath.Ext(file.Name()) != ".log" {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		// Remove files modified before the cutoff date
		if info.ModTime().Before(cutoffDate) {
			path := filepath.Join(h.baseDir, file.Name())
			if err := os.Remove(path); err != nil {
				fmt.Fprintf(os.Stderr, "error removing old log file %s: %v\n", path, err)
			}
		}
	}
}

// Close closes the current log file
func (h *DailyFileHandler) Close() error {
	h.fileMutex.Lock()
	defer h.fileMutex.Unlock()

	if h.currentFile != nil {
		err := h.currentFile.Close()
		h.currentFile = nil
		return err
	}
	return nil
}
