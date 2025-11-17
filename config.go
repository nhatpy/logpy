package logpy

import (
	"io"
	"os"
)

// OutputType defines where logs should be written
type OutputType string

const (
	OutputStdout OutputType = "stdout"
	OutputStderr OutputType = "stderr"
	OutputFile   OutputType = "file"
)

// FormatType defines the output format for logs
type FormatType string

const (
	FormatJSON    FormatType = "json"
	FormatConsole FormatType = "console"
)

// RotationMode defines how log files should be rotated
type RotationMode string

const (
	RotationSize  RotationMode = "size"  // Size-based rotation using lumberjack
	RotationDaily RotationMode = "daily" // Daily rotation based on date
)

// Config holds the configuration for creating a logger
type Config struct {
	// Level is the minimum log level to output
	Level Level

	// Format specifies the output format (json or console)
	Format FormatType

	// Output specifies where to write logs (stdout, stderr, or file)
	Output OutputType

	// OutputPath is the file path when Output is "file"
	OutputPath string

	// UseColor enables colored output for console format
	UseColor bool

	// ColorConfig allows customization of level colors
	ColorConfig ColorConfig

	// AddCaller includes caller information (file and line number)
	AddCaller bool

	// RotationMode specifies the rotation strategy: "size" or "daily"
	// Only used when Output is "file"
	RotationMode RotationMode

	// File rotation settings (used when Output is "file")
	// MaxSize is the maximum size in megabytes before rotation (for size-based rotation)
	MaxSize int

	// MaxBackups is the maximum number of old log files to retain (for size-based rotation)
	MaxBackups int

	// MaxAge is the maximum number of days to retain old log files
	MaxAge int

	// Compress determines if rotated files should be compressed (for size-based rotation)
	Compress bool

	// MultiOutput enables writing to both console and file
	MultiOutput bool
}

// DefaultConfig returns a configuration with sensible defaults
// Logs to BOTH console (with colors) and daily rotating files (with colors)
func DefaultConfig() Config {
	return Config{
		Level:        InfoLevel,
		Format:       FormatConsole,
		Output:       OutputFile,
		OutputPath:   "./logs",       // Just directory, no prefix
		UseColor:     true,            // Colors in both console and file
		ColorConfig:  DefaultColorConfig(),
		AddCaller:    true,
		RotationMode: RotationDaily,  // Daily rotation by default
		MaxSize:      100,             // 100 MB (for size-based rotation)
		MaxBackups:   3,               // Keep 3 old files (for size-based rotation)
		MaxAge:       28,              // Keep for 28 days
		Compress:     true,            // Compress old files (for size-based rotation)
		MultiOutput:  true,            // Log to BOTH console and file
	}
}

// DevelopmentConfig returns a configuration suitable for development
func DevelopmentConfig() Config {
	return Config{
		Level:       DebugLevel,
		Format:      FormatConsole,
		Output:      OutputStdout,
		UseColor:    true,
		ColorConfig: DefaultColorConfig(),
		AddCaller:   true,
		MaxSize:     100,
		MaxBackups:  3,
		MaxAge:      28,
		Compress:    true,
		MultiOutput: false,
	}
}

// ProductionConfig returns a configuration suitable for production
func ProductionConfig() Config {
	return Config{
		Level:       InfoLevel,
		Format:      FormatJSON,
		Output:      OutputFile,
		OutputPath:  "/var/log/app.log",
		UseColor:    false,
		ColorConfig: DefaultColorConfig(),
		AddCaller:   true,
		MaxSize:     100,
		MaxBackups:  5,
		MaxAge:      30,
		Compress:    true,
		MultiOutput: false,
	}
}

// isTerminal checks if stdout is a terminal
func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// getWriter returns the appropriate io.Writer based on config
func (c Config) getWriter() io.Writer {
	switch c.Output {
	case OutputStdout:
		return os.Stdout
	case OutputStderr:
		return os.Stderr
	case OutputFile:
		// Will be handled by FileHandler
		return nil
	default:
		return os.Stdout
	}
}
