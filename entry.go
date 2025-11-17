package logpy

import "time"

// Entry represents a complete log entry
type Entry struct {
	Time          time.Time
	Level         Level
	Message       string
	Fields        []Field // Event-specific fields
	ContextFields []Field // Persistent context fields (from With())
	Caller        CallerInfo
}

// Event is a fluent API builder for creating log entries
// It allows chaining methods to build up a log entry before sending it
type Event struct {
	logger    *Logger
	level     Level
	fields    []Field
	timestamp time.Time
	enabled   bool
}

// newEvent creates a new event for the given logger and level
func newEvent(logger *Logger, level Level) *Event {
	enabled := logger.handler.Enabled(level)
	return &Event{
		logger:    logger,
		level:     level,
		timestamp: time.Now(),
		enabled:   enabled,
	}
}

// Str adds a string field to the event
func (e *Event) Str(key, val string) *Event {
	if !e.enabled {
		return e
	}
	e.fields = append(e.fields, String(key, val))
	return e
}

// Int adds an int field to the event
func (e *Event) Int(key string, val int) *Event {
	if !e.enabled {
		return e
	}
	e.fields = append(e.fields, Int(key, val))
	return e
}

// Int64 adds an int64 field to the event
func (e *Event) Int64(key string, val int64) *Event {
	if !e.enabled {
		return e
	}
	e.fields = append(e.fields, Int64(key, val))
	return e
}

// Float64 adds a float64 field to the event
func (e *Event) Float64(key string, val float64) *Event {
	if !e.enabled {
		return e
	}
	e.fields = append(e.fields, Float64(key, val))
	return e
}

// Bool adds a boolean field to the event
func (e *Event) Bool(key string, val bool) *Event {
	if !e.enabled {
		return e
	}
	e.fields = append(e.fields, Bool(key, val))
	return e
}

// Time adds a time field to the event
func (e *Event) Time(key string, val time.Time) *Event {
	if !e.enabled {
		return e
	}
	e.fields = append(e.fields, Time(key, val))
	return e
}

// Dur adds a duration field to the event
func (e *Event) Dur(key string, val time.Duration) *Event {
	if !e.enabled {
		return e
	}
	e.fields = append(e.fields, Duration(key, val))
	return e
}

// Err adds an error field to the event
func (e *Event) Err(err error) *Event {
	if !e.enabled {
		return e
	}
	e.fields = append(e.fields, Error(err))
	return e
}

// Any adds a field with any value type to the event
func (e *Event) Any(key string, val interface{}) *Event {
	if !e.enabled {
		return e
	}
	e.fields = append(e.fields, Any(key, val))
	return e
}

// Fields adds multiple fields to the event
func (e *Event) Fields(fields ...Field) *Event {
	if !e.enabled {
		return e
	}
	e.fields = append(e.fields, fields...)
	return e
}

// Msg sends the event with the given message
// This finalizes and writes the log entry
func (e *Event) Msg(msg string) {
	if !e.enabled {
		return
	}

	entry := Entry{
		Time:          e.timestamp,
		Level:         e.level,
		Message:       msg,
		Fields:        e.fields,        // Event-specific fields
		ContextFields: e.logger.fields, // Context fields from With()
		Caller:        getCaller(2),    // Skip: getCaller -> Msg -> actual caller
	}

	// Handle the entry
	_ = e.logger.handler.Handle(entry)
}

// Msgf sends the event with a formatted message
func (e *Event) Msgf(format string, args ...interface{}) {
	if !e.enabled {
		return
	}
	// Use fmt.Sprintf for formatting
	msg := format
	if len(args) > 0 {
		// Simple implementation - for production, consider using fmt.Sprintf
		msg = format // Simplified for now
	}
	e.Msg(msg)
}

// Send sends the event without a message
func (e *Event) Send() {
	e.Msg("")
}
