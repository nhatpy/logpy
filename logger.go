package logpy

// Logger is the main logging interface
type Logger struct {
	handler Handler
	fields  []Field
}

// New creates a new logger with the provided handler
func New(handler Handler) *Logger {
	return &Logger{
		handler: handler,
		fields:  make([]Field, 0),
	}
}

// NewWithConfig creates a new logger with the provided configuration
func NewWithConfig(cfg Config) *Logger {
	var handler Handler

	switch cfg.Output {
	case OutputFile:
		// Check rotation mode
		if cfg.RotationMode == RotationDaily {
			// Daily rotation based on date
			baseDir := "./logs"
			filePrefix := "" // No prefix by default (just date.log)

			// Extract directory and optional prefix from OutputPath
			if cfg.OutputPath != "" {
				// If OutputPath ends with .log, it has a prefix
				if len(cfg.OutputPath) > 4 && cfg.OutputPath[len(cfg.OutputPath)-4:] == ".log" {
					// Extract directory and file prefix
					dir, file := splitPath(cfg.OutputPath)
					baseDir = dir
					// Remove .log extension to get prefix
					filePrefix = file[:len(file)-4]
				} else {
					// Just a directory path, no prefix
					baseDir = cfg.OutputPath
					filePrefix = "" // No prefix, just YYYY-MM-DD.log
				}
			}

			// Create daily file handler
			// File should have no colors if MultiOutput is enabled (colors go to console)
			// Otherwise, use the configured UseColor setting
			fileUseColor := cfg.UseColor && !cfg.MultiOutput
			dailyHandler, err := NewDailyFileHandler(
				baseDir,
				filePrefix,
				cfg.Level,
				cfg.MaxAge,
				fileUseColor,
				cfg.ColorConfig,
			)
			if err != nil {
				// Fallback to console handler on error
				handler = createConsoleHandler(cfg)
			} else {
				handler = dailyHandler
			}
		} else {
			// Size-based rotation using lumberjack
			handler = NewFileHandler(
				cfg.OutputPath,
				cfg.Level,
				cfg.MaxSize,
				cfg.MaxBackups,
				cfg.MaxAge,
				cfg.Compress,
			)
		}

		// If multi-output is enabled, also log to console
		if cfg.MultiOutput {
			// Console handler with colors enabled
			consoleHandler := NewConsoleHandlerWithConfig(cfg.Level, true, cfg.ColorConfig)
			handler = NewMultiHandler(handler, consoleHandler)
		}

	case OutputStdout, OutputStderr:
		if cfg.Format == FormatJSON {
			writer := cfg.getWriter()
			handler = NewJSONHandler(writer, cfg.Level)
		} else {
			handler = createConsoleHandler(cfg)
		}

	default:
		// Default to console handler
		handler = createConsoleHandler(cfg)
	}

	return &Logger{
		handler: handler,
		fields:  make([]Field, 0),
	}
}

// splitPath splits a file path into directory and filename
func splitPath(path string) (dir, file string) {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[:i], path[i+1:]
		}
	}
	return ".", path
}

// createConsoleHandler is a helper to create a console handler from config
func createConsoleHandler(cfg Config) Handler {
	return NewConsoleHandlerWithConfig(cfg.Level, cfg.UseColor, cfg.ColorConfig)
}

// Default creates a logger with default configuration
func Default() *Logger {
	return NewWithConfig(DefaultConfig())
}

// Development creates a logger with development configuration
func Development() *Logger {
	return NewWithConfig(DevelopmentConfig())
}

// Production creates a logger with production configuration
func Production() *Logger {
	return NewWithConfig(ProductionConfig())
}

// With creates a child logger with additional persistent fields
func (l *Logger) With(fields ...Field) *Logger {
	newFields := make([]Field, 0, len(l.fields)+len(fields))
	newFields = append(newFields, l.fields...)
	newFields = append(newFields, fields...)

	return &Logger{
		handler: l.handler,
		fields:  newFields,
	}
}

// Debug creates a debug level event
func (l *Logger) Debug() *Event {
	return newEvent(l, DebugLevel)
}

// Info creates an info level event
func (l *Logger) Info() *Event {
	return newEvent(l, InfoLevel)
}

// Warn creates a warn level event
func (l *Logger) Warn() *Event {
	return newEvent(l, WarnLevel)
}

// Error creates an error level event
func (l *Logger) Error() *Event {
	return newEvent(l, ErrorLevel)
}

// Global logger instance
var global = Default()

// SetGlobal sets the global logger instance
func SetGlobal(logger *Logger) {
	global = logger
}

// Global returns the global logger instance
func Global() *Logger {
	return global
}

// Log provides direct access to the global logger for quick logging
// Example: logpy.Log().Info().Str("key", "value").Msg("message")
func Log() *Logger {
	return global
}
