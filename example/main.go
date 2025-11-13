package main

import (
	"errors"
	"time"

	"github.com/nhatdo/logpy"
)

func main() {
	// Example 1: Using the default logger (Daily rotation to file with colors)
	println("=== Example 1: Default Logger (Daily Rotation with Colors) ===")
	defaultLogger := logpy.Default()
	defaultLogger.Info().Str("user", "john").Int("age", 30).Msg("User logged in")
	defaultLogger.Debug().Str("action", "login").Msg("Debug message")
	defaultLogger.Warn().Str("warning", "test").Msg("Warning message")
	defaultLogger.Error().Str("error", "test").Msg("Error message")
	println("Check ./logs/app-YYYY-MM-DD.log for colored output in file!")

	// Example 2: Development logger with colors
	println("\n=== Example 2: Development Logger (Colored Console) ===")
	devLogger := logpy.Development()
	devLogger.Debug().Str("component", "auth").Msg("Starting authentication process")
	devLogger.Info().Str("user", "alice").Int("user_id", 123).Msg("User authenticated successfully")
	devLogger.Warn().Str("ip", "192.168.1.1").Msg("Multiple failed login attempts detected")
	devLogger.Error().Err(errors.New("connection timeout")).Msg("Failed to connect to database")

	// Example 3: Custom configuration
	println("\n=== Example 3: Custom Configuration ===")
	customConfig := logpy.Config{
		Level:       logpy.DebugLevel,
		Format:      logpy.FormatConsole,
		Output:      logpy.OutputStdout,
		UseColor:    true,
		AddCaller:   true,
		ColorConfig: logpy.DefaultColorConfig(),
	}
	customLogger := logpy.NewWithConfig(customConfig)
	customLogger.Info().
		Str("service", "api-gateway").
		Int("port", 8080).
		Bool("tls_enabled", true).
		Msg("Server started")

	// Example 4: Child logger with persistent fields
	println("\n=== Example 4: Child Logger with Context ===")
	requestLogger := devLogger.With(
		logpy.String("request_id", "req-12345"),
		logpy.String("method", "GET"),
		logpy.String("path", "/api/users"),
	)
	requestLogger.Info().Msg("Request received")
	requestLogger.Info().Int("status", 200).Dur("latency", 45*time.Millisecond).Msg("Request completed")

	// Example 5: All field types
	println("\n=== Example 5: Different Field Types ===")
	devLogger.Info().
		Str("string_field", "hello").
		Int("int_field", 42).
		Int64("int64_field", 9876543210).
		Float64("float_field", 3.14159).
		Bool("bool_field", true).
		Time("time_field", time.Now()).
		Dur("duration_field", 2*time.Second).
		Any("any_field", map[string]string{"key": "value"}).
		Msg("Demonstrating all field types")

	// Example 6: File logger with rotation (commented to avoid creating files in example)
	println("\n=== Example 6: File Logger Configuration (not executed) ===")
	println("To use file logging with rotation:")
	println(`
	fileConfig := logpy.Config{
		Level:       logpy.InfoLevel,
		Format:      logpy.FormatJSON,
		Output:      logpy.OutputFile,
		OutputPath:  "/var/log/myapp/app.log",
		MaxSize:     100,  // 100 MB
		MaxBackups:  3,    // Keep 3 old files
		MaxAge:      28,   // Keep for 28 days
		Compress:    true, // Compress rotated files
		MultiOutput: true, // Also log to console
	}
	fileLogger := logpy.NewWithConfig(fileConfig)
	fileLogger.Info().Msg("This will be written to file with rotation")
	`)

	// Example 7: Global logger
	println("\n=== Example 7: Global Logger ===")
	logpy.SetGlobal(devLogger)
	logpy.Log().Info().Str("global", "true").Msg("Using global logger")

	// Example 8: Custom color configuration
	println("\n=== Example 8: Custom Colors ===")
	customColorConfig := logpy.Config{
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
	colorLogger := logpy.NewWithConfig(customColorConfig)
	colorLogger.Debug().Msg("Custom color for debug (Magenta)")
	colorLogger.Info().Msg("Custom color for info (Green)")
	colorLogger.Warn().Msg("Custom color for warn (Yellow)")
	colorLogger.Error().Msg("Custom color for error (Red)")

	// Example 9: Multiple handlers (console + file)
	println("\n=== Example 9: Multi-Output (Console + File) ===")
	println("This example logs to both console and file simultaneously")
	multiConfig := logpy.Config{
		Level:       logpy.InfoLevel,
		Format:      logpy.FormatJSON,
		Output:      logpy.OutputFile,
		OutputPath:  "./example.log",
		MaxSize:     10,
		MaxBackups:  2,
		MaxAge:      7,
		Compress:    false,
		MultiOutput: true, // Enable multi-output
	}
	multiLogger := logpy.NewWithConfig(multiConfig)
	multiLogger.Info().
		Str("feature", "multi-output").
		Msg("This message appears in both console and file!")
	println("Check example.log file for the output")
}
