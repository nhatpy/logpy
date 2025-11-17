# logpy - A Flexible and Feature-Rich Go Logging Library

`logpy` is a high-performance, structured logging library for Go, inspired by the best features of Zerolog, Zap, and Slog. It provides a clean fluent API, customizable output formats, and built-in daily file rotation support.

## Features

- **🎯 Log by Level**: Support for Debug, Info, Warn, and Error levels with filtering
- **🎨 Customizable Colors**: Full control over log level colors in console output
- **📅 Daily Log Rotation**: Automatic daily rotation with date-based filenames (e.g., `2025-11-17.log`)
- **📊 Structured Logging**: Easy-to-use fluent API for adding typed fields
- **🔄 Dual Output**: Log to console AND file simultaneously (default behavior)
- **📁 File Rotation**: Both daily rotation and size-based rotation supported
- **⚙️ Flexible Configuration**: Sensible defaults with full customization options
- **🚀 Zero Allocation**: Minimal allocations in hot paths for performance
- **📦 Minimal Dependencies**: Only requires `lumberjack` for size-based rotation

## Installation

```bash
go get github.com/nhatdo/logpy
```

## Quick Start

```go
package main

import "github.com/nhatdo/logpy"

func main() {
    // Use the default logger (console + file with daily rotation)
    logger := logpy.Default()
    logger.Info().Str("user", "john").Int("age", 30).Msg("User logged in")

    // Console: Colored output
    // File: ./logs/2025-11-17.log (plain text, no ANSI codes)
}
```

## Default Behavior ⭐

The default configuration logs to **BOTH console and file**:

- **Console**: Colored output (cyan timestamp, blue INFO, yellow WARN, red ERROR)
- **File**: Plain text at `./logs/2025-11-17.log` (no ANSI codes)
- **Rotation**: Daily (new file each day)
- **Cleanup**: Auto-delete logs older than 28 days

```go
logger := logpy.Default()
logger.Info().Msg("Logs to console AND file automatically!")
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

### 4. Daily File Rotation with Custom Prefix

```go
config := logpy.Config{
    Level:        logpy.InfoLevel,
    Format:       logpy.FormatConsole,
    Output:       logpy.OutputFile,
    OutputPath:   "./logs/myapp.log",   // Creates myapp-2025-11-17.log
    RotationMode: logpy.RotationDaily,  // Daily rotation
    UseColor:     false,                // Plain text in file
    MaxAge:       7,                    // Keep logs for 7 days
    MultiOutput:  false,                // File only
}

logger := logpy.NewWithConfig(config)
logger.Info().Msg("Logs to daily rotating file")
```

### 5. Size-Based Rotation (Traditional Lumberjack)

```go
config := logpy.Config{
    Level:        logpy.InfoLevel,
    Format:       logpy.FormatJSON,
    Output:       logpy.OutputFile,
    OutputPath:   "/var/log/myapp/app.log",
    RotationMode: logpy.RotationSize,  // Size-based rotation
    MaxSize:      100,  // 100 MB per file
    MaxBackups:   3,    // Keep 3 old files
    MaxAge:       28,   // Keep logs for 28 days
    Compress:     true, // Compress rotated files
}

logger := logpy.NewWithConfig(config)
logger.Info().Msg("Logs to size-rotated file")
```

### 6. Multi-Output (Console + File)

```go
config := logpy.Config{
    Level:        logpy.InfoLevel,
    Format:       logpy.FormatConsole,
    Output:       logpy.OutputFile,
    OutputPath:   "./logs",
    RotationMode: logpy.RotationDaily,
    UseColor:     true,
    MultiOutput:  true, // Enable both console and file output
}

logger := logpy.NewWithConfig(config)
logger.Info().Msg("This appears in both console and file")
// Console: With colors
// File: Plain text (no ANSI codes)
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
    Level       Level         // Minimum log level (DebugLevel, InfoLevel, WarnLevel, ErrorLevel)
    Format      FormatType    // Output format (FormatJSON, FormatConsole)
    Output      OutputType    // Output destination (OutputStdout, OutputStderr, OutputFile)
    OutputPath  string        // File path or directory (when Output is OutputFile)
    UseColor    bool          // Enable colored output (console format only)
    ColorConfig ColorConfig   // Custom color configuration
    AddCaller   bool          // Include caller information (file:line)

    // Rotation settings
    RotationMode RotationMode // "daily" or "size" rotation strategy
    MaxSize      int          // Maximum size in MB before rotation (size-based)
    MaxBackups   int          // Maximum number of old files to retain (size-based)
    MaxAge       int          // Maximum days to retain old files
    Compress     bool         // Compress rotated files with gzip (size-based)

    MultiOutput  bool         // Log to both console and file
}
```

### Preset Configurations

```go
// Default: Console + File (daily rotation), Info level
logger := logpy.Default()
// - Console: Colored output
// - File: ./logs/2025-11-17.log (plain text)
// - Daily rotation

// Development: Colored console output, Debug level
logger := logpy.Development()
// - Console only
// - All log levels (including Debug)
// - Colors enabled

// Production: JSON file output with size rotation, Info level
logger := logpy.Production()
// - File: /var/log/app.log
// - JSON format
// - Size-based rotation
```

## File Naming

### Daily Rotation

```go
// Simple (no prefix)
OutputPath: "./logs"
// Creates: ./logs/2025-11-17.log

// With prefix
OutputPath: "./logs/myapp.log"
// Creates: ./logs/myapp-2025-11-17.log
```

### Size-Based Rotation

```go
OutputPath: "/var/log/app.log"
// Creates: /var/log/app.log
// Rotated: /var/log/app.log.1, /var/log/app.log.2, etc.
```

## Color Behavior

| MultiOutput | Console | File | Use Case |
|-------------|---------|------|----------|
| `true` (default) | ✅ Colors | ❌ Plain text | Best for development/production |
| `false` + `UseColor=true` | ❌ | ✅ Colors | Terminal viewing of files |
| `false` + `UseColor=false` | ❌ | ❌ Plain text | Log aggregation systems |

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
ConsoleHandler / JSONHandler / DailyFileHandler / FileHandler / MultiHandler
    ↓
Formatter (JSON / Console)
    ↓
Writer (stdout / stderr / daily rotating file / size rotating file)
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
- Thread-safe file rotation

## Comparison with Other Libraries

| Feature | logpy | Zerolog | Zap | Slog |
|---------|-------|---------|-----|------|
| Fluent API | ✅ | ✅ | ❌ | ❌ |
| Typed Fields | ✅ | ✅ | ✅ | ✅ |
| Colors | ✅ (customizable) | ✅ (limited) | ✅ | ❌ (custom) |
| Daily Rotation | ✅ (built-in) | ❌ (external) | ❌ (external) | ❌ (external) |
| Multi-Output | ✅ (default) | ❌ (manual) | ❌ (manual) | ❌ (manual) |
| Dependencies | 1 (lumberjack) | 0 | 1 | 0 (stdlib) |
| Configuration | ✅ (struct-based) | Manual | Manual | Manual |

## Best Practices

1. **Use typed fields** instead of `Any()` for better performance
2. **Create child loggers** with `With()` for request/context scoping
3. **Use development config** in development for better readability
4. **Use production config** in production for structured JSON logs
5. **Set appropriate log levels** to control output verbosity
6. **Enable caller info** during development, consider disabling in production for performance
7. **Use daily rotation** for easier log management and analysis
8. **Keep MultiOutput enabled** (default) for best visibility during development

## Viewing Logs

```bash
# View plain text log
cat ./logs/2025-11-17.log

# Follow logs in real-time
tail -f ./logs/2025-11-17.log

# View logs with colors (if file has ANSI codes)
less -R ./logs/myapp-2025-11-17.log

# Search logs
grep "ERROR" ./logs/2025-11-17.log

# View specific date
cat ./logs/2025-11-16.log
```

## Example

See the [example](./example/main.go) directory for a complete working example demonstrating all features.

Run it:
```bash
cd example
go run main.go
```

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

## Changelog

### v1.1.0 (Latest)
- ✨ Added daily log rotation with date-based filenames
- ✨ Multi-output (console + file) as default behavior
- ✨ Smart color handling: colors in console, plain text in files
- 🎨 Improved default configuration for better developer experience
- 📝 Comprehensive examples and documentation

### v1.0.0
- Initial release with basic logging features
