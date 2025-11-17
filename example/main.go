package main

import (
	"errors"
	"time"

	"github.com/nhatdo/logpy"
)

func main() {
	// Example 1: Default Logger - Daily Rotation with Colors in File
	println("=== Example 1: Default Logger (Daily Rotation with Colors) ===")
	println("Creates: ./logs/app-2025-11-07.log (with ANSI color codes)")
	defaultLogger := logpy.Default()
	defaultLogger.Info().Str("user", "john").Int("age", 30).Msg("User logged in")
	defaultLogger.Warn().Str("warning", "test").Msg("Warning message")
	defaultLogger.Error().Str("error", "test").Msg("Error message")
	println("✓ Logs written to ./logs/app-" + time.Now().Format("2006-01-02") + ".log")
	println("✓ New file created automatically each day")
	println("✓ Old logs auto-deleted after 28 days")

	// Example 2: Development Logger - Console Output with Colors
	println("\n=== Example 2: Development Logger (Colored Console) ===")
	devLogger := logpy.Development()
	devLogger.Debug().Str("component", "auth").Msg("Starting authentication process")
	devLogger.Info().Str("user", "alice").Int("user_id", 123).Msg("User authenticated successfully")
	devLogger.Warn().Str("ip", "192.168.1.1").Msg("Multiple failed login attempts detected")
	devLogger.Error().Err(errors.New("connection timeout")).Msg("Failed to connect to database")
	println("✓ Outputs to console (stdout) with colors")
	println("✓ DEBUG level enabled")

	// Example 3: Daily Rotation with Custom Filename Prefix
	println("\n=== Example 3: Daily Rotation - Custom Filename ===")
	println("Creates: ./logs/myservice-2025-11-07.log")
	customDailyConfig := logpy.Config{
		Level:        logpy.InfoLevel,
		Format:       logpy.FormatConsole,
		Output:       logpy.OutputFile,
		OutputPath:   "./logs/myservice.log", // Prefix becomes "myservice"
		RotationMode: logpy.RotationDaily,
		UseColor:     true,
		MaxAge:       7, // Keep logs for 7 days
		AddCaller:    true,
		ColorConfig:  logpy.DefaultColorConfig(),
	}
	customLogger := logpy.NewWithConfig(customDailyConfig)
	customLogger.Info().
		Str("service", "api-gateway").
		Int("port", 8080).
		Bool("tls_enabled", true).
		Msg("Server started")
	println("✓ Custom prefix: myservice-2025-11-07.log")
	println("✓ Logs retained for 7 days only")

	// Example 4: Size-Based Rotation (Traditional Lumberjack)
	println("\n=== Example 4: Size-Based Rotation (not daily) ===")
	println("Creates: /tmp/sizerotate.log (rotates at 10MB)")
	sizeRotateConfig := logpy.Config{
		Level:        logpy.InfoLevel,
		Format:       logpy.FormatJSON,
		Output:       logpy.OutputFile,
		OutputPath:   "/tmp/sizerotate.log",
		RotationMode: logpy.RotationSize, // Size-based, not daily
		MaxSize:      10,                  // 10 MB per file
		MaxBackups:   3,                   // Keep 3 old files
		MaxAge:       30,                  // Keep for 30 days
		Compress:     true,                // Compress old files
		UseColor:     false,               // No colors for JSON
	}
	sizeLogger := logpy.NewWithConfig(sizeRotateConfig)
	sizeLogger.Info().Str("mode", "size-based").Msg("Using lumberjack rotation")
	println("✓ Rotates when file reaches 10MB")
	println("✓ Keeps 3 backup files")
	println("✓ JSON format without colors")

	// Example 5: Child Logger with Persistent Fields
	println("\n=== Example 5: Child Logger with Context ===")
	requestLogger := devLogger.With(
		logpy.String("request_id", "req-12345"),
		logpy.String("method", "GET"),
		logpy.String("path", "/api/users"),
	)
	requestLogger.Info().Msg("Request received")
	requestLogger.Info().Int("status", 200).Dur("latency", 45*time.Millisecond).Msg("Request completed")
	println("✓ All logs include: request_id, method, path")

	// Example 6: All Field Types
	println("\n=== Example 6: Different Field Types ===")
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
	println("✓ Supports: Str, Int, Int64, Float64, Bool, Time, Dur, Any, Err")

	// Example 7: Multi-Output (Console + File Simultaneously)
	println("\n=== Example 7: Multi-Output (Console + File) ===")
	println("Logs to BOTH console AND daily file simultaneously")
	dualConfig := logpy.Config{
		Level:        logpy.InfoLevel,
		Format:       logpy.FormatConsole,
		Output:       logpy.OutputFile,
		OutputPath:   "./logs/dual.log",
		RotationMode: logpy.RotationDaily,
		UseColor:     true,
		MultiOutput:  true, // Enable dual output
		MaxAge:       14,
	}
	dualLogger := logpy.NewWithConfig(dualConfig)
	dualLogger.Info().
		Str("feature", "multi-output").
		Msg("This appears in BOTH console and file!")
	println("✓ Console: colored output to stdout")
	println("✓ File: ./logs/dual-2025-11-07.log (with colors)")

	// Example 8: Global Logger
	println("\n=== Example 8: Global Logger ===")
	logpy.SetGlobal(devLogger)
	logpy.Log().Info().Str("global", "true").Msg("Using global logger")
	println("✓ Use logpy.Log() anywhere in your application")

	// Example 9: Custom Color Configuration
	println("\n=== Example 9: Custom Colors ===")
	customColorConfig := logpy.Config{
		Level:    logpy.DebugLevel,
		Format:   logpy.FormatConsole,
		Output:   logpy.OutputStdout,
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
	println("✓ Fully customizable color scheme")

	// Example 10: Daily Rotation without Colors (for plain text)
	println("\n=== Example 10: Daily Rotation without Colors ===")
	println("Creates: ./logs/plain-2025-11-07.log (no ANSI codes)")
	plainConfig := logpy.Config{
		Level:        logpy.InfoLevel,
		Format:       logpy.FormatConsole,
		Output:       logpy.OutputFile,
		OutputPath:   "./logs/plain.log",
		RotationMode: logpy.RotationDaily,
		UseColor:     false, // No colors
		MaxAge:       30,
	}
	plainLogger := logpy.NewWithConfig(plainConfig)
	plainLogger.Info().Str("format", "plain").Msg("No color codes in this file")
	println("✓ Plain text format for editors that don't support ANSI")

	// Summary
	println("\n=== Summary ===")
	println("Log files created:")
	println("  - ./logs/app-" + time.Now().Format("2006-01-02") + ".log (default, with colors)")
	println("  - ./logs/myservice-" + time.Now().Format("2006-01-02") + ".log (custom prefix)")
	println("  - ./logs/dual-" + time.Now().Format("2006-01-02") + ".log (multi-output)")
	println("  - ./logs/plain-" + time.Now().Format("2006-01-02") + ".log (no colors)")
	println("  - /tmp/sizerotate.log (size-based rotation)")
	println("\nView logs with colors: cat ./logs/app-" + time.Now().Format("2006-01-02") + ".log")
	println("Follow logs in real-time: tail -f ./logs/app-" + time.Now().Format("2006-01-02") + ".log")
}
