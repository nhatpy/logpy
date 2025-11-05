# logpy - A Flexible and Feature-Rich Go Logging Library

`logpy` is a high-performance, structured logging library for Go, inspired by the best features of Zerolog, Zap, and Slog. It provides a clean fluent API, customizable output formats, and built-in file rotation support.

## Features

- **Log by Level**: Support for Debug, Info, Warn, and Error levels with filtering
- **Customizable Colors**: Full control over log level colors in console output
- **Detailed Context**: Automatic timestamp and caller information (file and line number)
- **Structured Logging**: Easy-to-use fluent API for adding typed fields
- **File Rotation**: Built-in support for size-based and time-based log rotation
- **Flexible Configuration**: Sensible defaults with full customization options
- **Multiple Outputs**: Log to console and file simultaneously
- **Zero Dependencies**: Only requires `lumberjack` for file rotation (optional)

## Installation

```bash
go get github.com/nhatdo/logpy
```

## Quick Start

```go
package main

import "github.com/nhatdo/logpy"

func main() {
    // Use the default logger (JSON output to stdout)
    logger := logpy.Default()
    logger.Info().Str("user", "john").Int("age", 30).Msg("User logged in")

    // Or use the development logger (colored console output)
    devLogger := logpy.Development()
    devLogger.Debug().Str("component", "auth").Msg("Starting authentication")
}
```

## Usage Examples

### 1. Basic Logging with Different Levels

```go
logger := logpy.Development()

logger.Debug().Msg("Debug message")
logger.Info().Msg("Info message")
logger.Warn().Msg("Warning message")
logger.Error().Msg("Error message")
```

### 2. Structured Logging (Adding Fields)

```go
logger.Info().
    Str("service", "api-gateway").
    Int("port", 8080).
    Bool("tls_enabled", true).
    Dur("startup_time", 2*time.Second).
    Msg("Server started")
```

### 3. Child Logger with Persistent Fields

```go
// Create a child logger with request context
requestLogger := logger.With(
    logpy.String("request_id", "req-12345"),
    logpy.String("method", "GET"),
    logpy.String("path", "/api/users"),
)

// All logs from this logger will include the request context
requestLogger.Info().Msg("Request received")
requestLogger.Info().Int("status", 200).Msg("Request completed")
```

### 4. Custom Configuration

```go
config := logpy.Config{
    Level:       logpy.InfoLevel,
    Format:      logpy.FormatConsole,
    Output:      logpy.OutputStdout,
    UseColor:    true,
    AddCaller:   true,
    ColorConfig: logpy.DefaultColorConfig(),
}

logger := logpy.NewWithConfig(config)
```

### 5. File Logging with Rotation

```go
config := logpy.Config{
    Level:       logpy.InfoLevel,
    Format:      logpy.FormatJSON,
    Output:      logpy.OutputFile,
    OutputPath:  "/var/log/myapp/app.log",
    MaxSize:     100,  // 100 MB per file
    MaxBackups:  3,    // Keep 3 old files
    MaxAge:      28,   // Keep logs for 28 days
    Compress:    true, // Compress rotated files
}

logger := logpy.NewWithConfig(config)
logger.Info().Msg("This will be written to a rotating file")
```

### 6. Multi-Output (Console + File)

```go
config := logpy.Config{
    Level:       logpy.InfoLevel,
    Format:      logpy.FormatJSON,
    Output:      logpy.OutputFile,
    OutputPath:  "/var/log/app.log",
    MultiOutput: true, // Enable both console and file output
}

logger := logpy.NewWithConfig(config)
logger.Info().Msg("This appears in both console and file")
```

### 7. Custom Colors

```go
config := logpy.Config{
    Level:  logpy.DebugLevel,
    Format: logpy.FormatConsole,
    Output: logpy.OutputStdout,
    UseColor: true,
    ColorConfig: logpy.ColorConfig{
        Debug: "\033[35m", // Magenta
        Info:  "\033[32m", // Green
        Warn:  "\033[33m", // Yellow
        Error: "\033[31m", // Red
        Reset: "\033[0m",
    },
}

logger := logpy.NewWithConfig(config)
```

### 8. Global Logger

```go
// Set a global logger
logpy.SetGlobal(logpy.Development())

// Use it anywhere in your application
logpy.Log().Info().Str("global", "true").Msg("Using global logger")
```

### 9. All Field Types

```go
logger.Info().
    Str("string", "hello").
    Int("int", 42).
    Int64("int64", 9876543210).
    Float64("float", 3.14159).
    Bool("bool", true).
    Time("time", time.Now()).
    Dur("duration", 2*time.Second).
    Err(errors.New("some error")).
    Any("custom", map[string]string{"key": "value"}).
    Msg("All field types")
```

## Configuration Options

### Config Struct

```go
type Config struct {
    Level       Level      // Minimum log level (DebugLevel, InfoLevel, WarnLevel, ErrorLevel)
    Format      FormatType // Output format (FormatJSON, FormatConsole)
    Output      OutputType // Output destination (OutputStdout, OutputStderr, OutputFile)
    OutputPath  string     // File path (when Output is OutputFile)
    UseColor    bool       // Enable colored output (console format only)
    ColorConfig ColorConfig // Custom color configuration
    AddCaller   bool       // Include caller information (file:line)

    // File rotation settings (when Output is OutputFile)
    MaxSize     int  // Maximum size in MB before rotation
    MaxBackups  int  // Maximum number of old files to retain
    MaxAge      int  // Maximum number of days to retain old files
    Compress    bool // Compress rotated files with gzip

    MultiOutput bool // Log to both console and file
}
```

### Preset Configurations

```go
// Default: JSON output to stdout, Info level
logger := logpy.Default()

// Development: Colored console output, Debug level
logger := logpy.Development()

// Production: JSON file output with rotation, Info level
logger := logpy.Production()
```

## API Reference

### Logger Methods

- `Debug()` - Create a debug level event
- `Info()` - Create an info level event
- `Warn()` - Create a warn level event
- `Error()` - Create an error level event
- `With(fields ...Field)` - Create a child logger with persistent fields

### Event Methods (Chainable)

- `Str(key, val string)` - Add a string field
- `Int(key string, val int)` - Add an int field
- `Int64(key string, val int64)` - Add an int64 field
- `Float64(key string, val float64)` - Add a float64 field
- `Bool(key string, val bool)` - Add a boolean field
- `Time(key string, val time.Time)` - Add a time field
- `Dur(key string, val time.Duration)` - Add a duration field
- `Err(err error)` - Add an error field
- `Any(key string, val interface{})` - Add any value (uses reflection)
- `Msg(msg string)` - Send the event with a message
- `Send()` - Send the event without a message

### Field Constructors

```go
logpy.String(key, val string)
logpy.Int(key string, val int)
logpy.Int64(key string, val int64)
logpy.Float64(key string, val float64)
logpy.Bool(key string, val bool)
logpy.Time(key string, val time.Time)
logpy.Duration(key string, val time.Duration)
logpy.Error(err error)
logpy.Any(key string, val interface{})
```

## Architecture

`logpy` uses a handler-based architecture inspired by Go's `slog`:

```
Logger (Frontend)
    ↓
Handler Interface (Backend)
    ↓
ConsoleHandler / JSONHandler / FileHandler / MultiHandler
    ↓
Formatter (JSON / Console)
    ↓
Writer (stdout / stderr / file with rotation)
```

This design provides:
- Clean separation of concerns
- Easy extensibility (custom handlers and formatters)
- Multiple output support
- Testability

## Performance

`logpy` is designed for performance with:
- Minimal allocations in hot paths
- Early level filtering (disabled events are no-ops)
- Efficient field builders
- Optional caller detection

## Comparison with Other Libraries

| Feature | logpy | Zerolog | Zap | Slog |
|---------|-------|---------|-----|------|
| Fluent API | ✅ | ✅ | ❌ | ❌ |
| Typed Fields | ✅ | ✅ | ✅ | ✅ |
| Colors | ✅ (customizable) | ✅ (limited) | ✅ | ❌ (custom) |
| File Rotation | ✅ (built-in) | ❌ (external) | ❌ (external) | ❌ (external) |
| Dependencies | 1 (lumberjack) | 0 | 1 | 0 (stdlib) |
| Configuration | ✅ (struct-based) | Manual | Manual | Manual |

## Best Practices

1. **Use typed fields** instead of `Any()` for better performance
2. **Create child loggers** with `With()` for request/context scoping
3. **Use development config** in development for better readability
4. **Use production config** in production for structured JSON logs
5. **Set appropriate log levels** to control output verbosity
6. **Enable caller info** during development, consider disabling in production for performance

## Example

See the [example](./example/main.go) directory for a complete working example demonstrating all features.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

MIT License

## Acknowledgments

`logpy` is inspired by:
- [Zerolog](https://github.com/rs/zerolog) - Fluent API and zero-allocation design
- [Zap](https://github.com/uber-go/zap) - Strongly-typed fields and performance focus
- [Slog](https://pkg.go.dev/log/slog) - Handler-based architecture and standard library approach
- [Lumberjack](https://github.com/natefinch/lumberjack) - File rotation implementation
