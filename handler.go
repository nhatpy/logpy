package logpy

import (
	"io"
	"os"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Handler is an interface for processing log entries
type Handler interface {
	// Enabled reports whether the handler handles records at the given level
	Enabled(level Level) bool
	// Handle processes a log entry
	Handle(entry Entry) error
	// WithFields returns a new handler with additional persistent fields
	WithFields(fields []Field) Handler
}

// baseHandler provides common functionality for all handlers
type baseHandler struct {
	level     Level
	formatter Formatter
	writer    io.Writer
	mu        sync.Mutex
}

// Enabled implements the Handler interface
func (h *baseHandler) Enabled(level Level) bool {
	return level >= h.level
}

// Handle implements the Handler interface
func (h *baseHandler) Handle(entry Entry) error {
	if !h.Enabled(entry.Level) {
		return nil
	}

	// Format the entry
	data, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}

	// Write to output (thread-safe)
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err = h.writer.Write(data)
	return err
}

// WithFields implements the Handler interface
func (h *baseHandler) WithFields(fields []Field) Handler {
	// For base handler, we don't modify the handler itself
	// The fields will be managed by the logger
	return h
}

// ConsoleHandler is a handler that writes to console with optional colors
type ConsoleHandler struct {
	*baseHandler
}

// NewConsoleHandler creates a new console handler
func NewConsoleHandler(level Level, useColor bool) *ConsoleHandler {
	formatter := &ConsoleFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		AddCaller:       true,
		UseColor:        useColor,
		ColorConfig:     DefaultColorConfig(),
	}

	return &ConsoleHandler{
		baseHandler: &baseHandler{
			level:     level,
			formatter: formatter,
			writer:    os.Stdout,
		},
	}
}

// NewConsoleHandlerWithConfig creates a console handler with custom configuration
func NewConsoleHandlerWithConfig(level Level, useColor bool, colorConfig ColorConfig) *ConsoleHandler {
	formatter := &ConsoleFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		AddCaller:       true,
		UseColor:        useColor,
		ColorConfig:     colorConfig,
	}

	return &ConsoleHandler{
		baseHandler: &baseHandler{
			level:     level,
			formatter: formatter,
			writer:    os.Stdout,
		},
	}
}

// JSONHandler is a handler that writes JSON formatted logs
type JSONHandler struct {
	*baseHandler
}

// NewJSONHandler creates a new JSON handler that writes to the specified writer
func NewJSONHandler(writer io.Writer, level Level) *JSONHandler {
	formatter := &JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00", // ISO 8601
		AddCaller:       true,
	}

	return &JSONHandler{
		baseHandler: &baseHandler{
			level:     level,
			formatter: formatter,
			writer:    writer,
		},
	}
}

// FileHandler is a handler that writes to a file with rotation support
type FileHandler struct {
	*baseHandler
	rotator *lumberjack.Logger
}

// NewFileHandler creates a new file handler with rotation support
func NewFileHandler(filename string, level Level, maxSize, maxBackups, maxAge int, compress bool) *FileHandler {
	rotator := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,    // MB
		MaxBackups: maxBackups, // Number of old files to keep
		MaxAge:     maxAge,     // Days to retain old files
		Compress:   compress,   // Compress rotated files
		LocalTime:  true,       // Use local time for filenames
	}

	formatter := &JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		AddCaller:       true,
	}

	return &FileHandler{
		baseHandler: &baseHandler{
			level:     level,
			formatter: formatter,
			writer:    rotator,
		},
		rotator: rotator,
	}
}

// Close closes the file handler and flushes any buffered data
func (h *FileHandler) Close() error {
	return h.rotator.Close()
}

// MultiHandler sends log entries to multiple handlers
type MultiHandler struct {
	handlers []Handler
	level    Level
}

// NewMultiHandler creates a handler that writes to multiple handlers
func NewMultiHandler(handlers ...Handler) *MultiHandler {
	// Find the minimum level among all handlers
	minLevel := ErrorLevel
	for _, h := range handlers {
		for l := DebugLevel; l <= ErrorLevel; l++ {
			if h.Enabled(l) {
				if l < minLevel {
					minLevel = l
				}
				break
			}
		}
	}

	return &MultiHandler{
		handlers: handlers,
		level:    minLevel,
	}
}

// Enabled implements the Handler interface
func (h *MultiHandler) Enabled(level Level) bool {
	// Return true if ANY handler is enabled at this level
	for _, handler := range h.handlers {
		if handler.Enabled(level) {
			return true
		}
	}
	return false
}

// Handle implements the Handler interface
func (h *MultiHandler) Handle(entry Entry) error {
	var lastErr error
	for _, handler := range h.handlers {
		if err := handler.Handle(entry); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// WithFields implements the Handler interface
func (h *MultiHandler) WithFields(fields []Field) Handler {
	newHandlers := make([]Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithFields(fields)
	}
	return NewMultiHandler(newHandlers...)
}
