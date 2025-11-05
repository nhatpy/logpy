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
		// File output with rotation
		handler = NewFileHandler(
			cfg.OutputPath,
			cfg.Level,
			cfg.MaxSize,
			cfg.MaxBackups,
			cfg.MaxAge,
			cfg.Compress,
		)

		// If multi-output is enabled, also log to console
		if cfg.MultiOutput {
			consoleHandler := createConsoleHandler(cfg)
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
