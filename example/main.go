package main

import (
	"errors"
	"time"

	"github.com/nhatdo/logpy"
)

func main() {
	// Example 1: Default Logger - BOTH Console (with colors) AND File (plain text)
	println("=== Example 1: Default Logger (Console + File) ===")
	println("Default config logs to BOTH console and file automatically!")
	defaultLogger := logpy.Default()
	defaultLogger.Info().Str("user", "john").Int("age", 30).Msg("User logged in")
	defaultLogger.Warn().Str("warning", "test").Msg("Warning message")
	defaultLogger.Error().Str("error", "database").Msg("Connection failed")
	println("✓ Console: colored output (cyan timestamp, blue INFO, yellow WARN, red ERROR)")
	println("✓ File: ./logs/" + time.Now().Format("2006-01-02") + ".log (plain text, no ANSI codes)")
	println("✓ New file created each day automatically")
	println("✓ Old logs deleted after 28 days")

	// Example 2: Development Logger - Console Only with Colors
	println("\n=== Example 2: Development Logger (Console Only) ===")
	devLogger := logpy.Development()
	devLogger.Debug().Str("component", "auth").Msg("Starting authentication process")
	devLogger.Info().Str("user", "alice").Int("user_id", 123).Msg("User authenticated successfully")
	devLogger.Warn().Str("ip", "192.168.1.1").Msg("Multiple failed login attempts detected")
	devLogger.Error().Err(errors.New("connection timeout")).Msg("Failed to connect to database")
	println("✓ Console output only (stdout) with colors")
	println("✓ DEBUG level enabled")

	// Example 3: Daily File Rotation with Custom Prefix
	println("\n=== Example 3: File Only - Daily Rotation with Custom Prefix ===")
	customDailyConfig := logpy.Config{
		Level:        logpy.InfoLevel,
		Format:       logpy.FormatConsole,
		Output:       logpy.OutputFile,
		OutputPath:   "./logs/myservice.log", // Creates myservice-2025-11-17.log
		RotationMode: logpy.RotationDaily,
		UseColor:     true,                    // Colors in file since MultiOutput=false
		MaxAge:       7,                       // Keep logs for 7 days
		AddCaller:    true,
		ColorConfig:  logpy.DefaultColorConfig(),
		MultiOutput:  false,                   // File only
	}
	customLogger := logpy.NewWithConfig(customDailyConfig)
	customLogger.Info().
		Str("service", "api-gateway").
		Int("port", 8080).
		Bool("tls_enabled", true).
		Msg("Server started")
	println("✓ File: ./logs/myservice-" + time.Now().Format("2006-01-02") + ".log")
	println("✓ WITH colors in file (since MultiOutput=false)")

	// Example 4: Size-Based Rotation (JSON format)
	println("\n=== Example 4: Size-Based Rotation (Traditional) ===")
	sizeRotateConfig := logpy.Config{
		Level:        logpy.InfoLevel,
		Format:       logpy.FormatJSON,
		Output:       logpy.OutputFile,
		OutputPath:   "/tmp/sizerotate.log",
		RotationMode: logpy.RotationSize, // Size-based
		MaxSize:      10,                  // 10 MB
		MaxBackups:   3,
		MaxAge:       30,
		Compress:     true,
		UseColor:     false,
		MultiOutput:  false,
	}
	sizeLogger := logpy.NewWithConfig(sizeRotateConfig)
	sizeLogger.Info().Str("mode", "size-based").Msg("Using lumberjack rotation")
	println("✓ JSON format, rotates at 10MB")

	// Example 5: Child Logger with Context
	println("\n=== Example 5: Child Logger with Persistent Fields ===")
	requestLogger := devLogger.With(
		logpy.String("request_id", "req-12345"),
		logpy.String("method", "GET"),
		logpy.String("path", "/api/users"),
	)
	requestLogger.Info().Msg("Request received")
	requestLogger.Info().Int("status", 200).Dur("latency", 45*time.Millisecond).Msg("Request completed")
	println("✓ All logs include: request_id, method, path")

	// Example 6: All Field Types
	println("\n=== Example 6: All Field Types ===")
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
	println("✓ Field types: Str, Int, Int64, Float64, Bool, Time, Dur, Any, Err")

	// Example 7: Global Logger
	println("\n=== Example 7: Global Logger ===")
	logpy.SetGlobal(devLogger)
	logpy.Log().Info().Str("global", "true").Msg("Using global logger")
	println("✓ Use logpy.Log() anywhere in your application")

	// Example 8: Custom Colors
	println("\n=== Example 8: Custom Color Configuration ===")
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
	colorLogger.Debug().Msg("Custom debug color (Magenta)")
	colorLogger.Info().Msg("Custom info color (Green)")
	colorLogger.Warn().Msg("Custom warn color (Yellow)")
	colorLogger.Error().Msg("Custom error color (Red)")
	println("✓ Fully customizable color scheme")

	// Summary
	println("\n=== Summary ===")
	println("Key Points:")
	println("1. DEFAULT CONFIG = Console (colors) + File (plain text)")
	println("2. File naming: YYYY-MM-DD.log (e.g., " + time.Now().Format("2006-01-02") + ".log)")
	println("3. Daily rotation: New file each day")
	println("4. No ANSI codes in files when MultiOutput=true")
	println("5. Colors appear in files only when MultiOutput=false and UseColor=true")
	println("\nLog files created:")
	println("  - ./logs/" + time.Now().Format("2006-01-02") + ".log (default, plain text)")
	println("  - ./logs/myservice-" + time.Now().Format("2006-01-02") + ".log (with colors)")
	println("  - /tmp/sizerotate.log (JSON, size-based)")
	println("\nView logs:")
	println("  cat ./logs/" + time.Now().Format("2006-01-02") + ".log")
	println("  tail -f ./logs/" + time.Now().Format("2006-01-02") + ".log")
}
