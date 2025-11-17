package logpy

import (
	"encoding/json"
	"fmt"
	"time"
)

// Color codes for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[37m"
	colorCyan   = "\033[36m"
)

// ColorConfig allows customization of log level colors
type ColorConfig struct {
	Debug string
	Info  string
	Warn  string
	Error string
	Reset string
}

// DefaultColorConfig returns the default color configuration
func DefaultColorConfig() ColorConfig {
	return ColorConfig{
		Debug: colorGray,
		Info:  colorBlue,
		Warn:  colorYellow,
		Error: colorRed,
		Reset: colorReset,
	}
}

// Formatter is an interface for formatting log entries
type Formatter interface {
	Format(entry Entry) ([]byte, error)
}

// JSONFormatter formats log entries as JSON
type JSONFormatter struct {
	TimestampFormat string
	AddCaller       bool
}

// Format implements the Formatter interface for JSON output
func (f *JSONFormatter) Format(entry Entry) ([]byte, error) {
	m := make(map[string]interface{})

	// Add timestamp
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}
	m["timestamp"] = entry.Time.Format(timestampFormat)

	// Add level
	m["level"] = entry.Level.String()

	// Add message
	if entry.Message != "" {
		m["message"] = entry.Message
	}

	// Add caller info
	if f.AddCaller {
		m["caller"] = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
	}

	// Add event-specific fields
	for _, field := range entry.Fields {
		m[field.Key] = field.Value
	}

	// Add context fields under "context" key
	if len(entry.ContextFields) > 0 {
		contextData := make(map[string]interface{})
		for _, field := range entry.ContextFields {
			contextData[field.Key] = field.Value
		}
		m["context"] = contextData
	}

	// Marshal to JSON
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	// Add newline
	data = append(data, '\n')
	return data, nil
}

// ConsoleFormatter formats log entries for console output with colors
type ConsoleFormatter struct {
	TimestampFormat string
	AddCaller       bool
	UseColor        bool
	ColorConfig     ColorConfig
}

// Format implements the Formatter interface for console output
func (f *ConsoleFormatter) Format(entry Entry) ([]byte, error) {
	var output string

	// Get color for level
	levelColor := ""
	if f.UseColor {
		switch entry.Level {
		case DebugLevel:
			levelColor = f.ColorConfig.Debug
		case InfoLevel:
			levelColor = f.ColorConfig.Info
		case WarnLevel:
			levelColor = f.ColorConfig.Warn
		case ErrorLevel:
			levelColor = f.ColorConfig.Error
		}
	}

	// Format timestamp
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = "2006-01-02 15:04:05"
	}
	timestamp := entry.Time.Format(timestampFormat)

	// Build output string
	if f.UseColor {
		output = fmt.Sprintf("%s[%s] %s%-5s%s", colorCyan, timestamp, levelColor, entry.Level.String(), f.ColorConfig.Reset)
	} else {
		output = fmt.Sprintf("[%s] %-5s", timestamp, entry.Level.String())
	}

	// Add caller info
	if f.AddCaller {
		output += fmt.Sprintf(" %s:%d", entry.Caller.File, entry.Caller.Line)
	}

	// Add message
	if entry.Message != "" {
		output += " " + entry.Message
	}

	// Add event-specific fields first
	if len(entry.Fields) > 0 {
		for _, field := range entry.Fields {
			output += fmt.Sprintf(" %s=%v", field.Key, field.Value)
		}
	}

	// Add context fields (separated with | symbol)
	if len(entry.ContextFields) > 0 {
		output += " |"
		for _, field := range entry.ContextFields {
			output += fmt.Sprintf(" %s=%v", field.Key, field.Value)
		}
	}

	output += "\n"
	return []byte(output), nil
}
