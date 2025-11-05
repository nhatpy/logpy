package logpy

import "time"

// FieldType represents the type of a field value
type FieldType uint8

const (
	StringType FieldType = iota
	IntType
	Int64Type
	Float64Type
	BoolType
	TimeType
	DurationType
	ErrorType
	AnyType
)

// Field represents a strongly-typed key-value pair for structured logging
type Field struct {
	Key   string
	Type  FieldType
	Value interface{}
}

// String creates a string field
func String(key, val string) Field {
	return Field{Key: key, Type: StringType, Value: val}
}

// Int creates an int field
func Int(key string, val int) Field {
	return Field{Key: key, Type: IntType, Value: val}
}

// Int64 creates an int64 field
func Int64(key string, val int64) Field {
	return Field{Key: key, Type: Int64Type, Value: val}
}

// Float64 creates a float64 field
func Float64(key string, val float64) Field {
	return Field{Key: key, Type: Float64Type, Value: val}
}

// Bool creates a boolean field
func Bool(key string, val bool) Field {
	return Field{Key: key, Type: BoolType, Value: val}
}

// Time creates a time field
func Time(key string, val time.Time) Field {
	return Field{Key: key, Type: TimeType, Value: val}
}

// Duration creates a duration field
func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Type: DurationType, Value: val}
}

// Error creates an error field
func Error(err error) Field {
	if err == nil {
		return Field{Key: "error", Type: ErrorType, Value: nil}
	}
	return Field{Key: "error", Type: ErrorType, Value: err.Error()}
}

// Any creates a field with any value type (uses reflection, slower)
func Any(key string, val interface{}) Field {
	return Field{Key: key, Type: AnyType, Value: val}
}
