package main

import (
	"errors"
	"time"

	"github.com/nhatdo/logpy"
)

func main() {
	println("╔════════════════════════════════════════════════════════════════╗")
	println("║          LogPy - Comprehensive Logging Examples               ║")
	println("╚════════════════════════════════════════════════════════════════╝")

	// Example 1: Default Logger - BOTH Console AND File
	println("\n=== Example 1: Default Logger (Console + File) ===")
	println("The default configuration logs to BOTH console and file!")
	defaultLogger := logpy.Default()
	defaultLogger.Info().Str("user", "john").Int("age", 30).Msg("User logged in")
	defaultLogger.Warn().Str("warning", "high_memory").Int("usage_mb", 512).Msg("Memory usage warning")
	defaultLogger.Error().Str("error", "database").Msg("Connection failed")

	println("\n✓ Console: Colored output")
	println("  - Cyan timestamp")
	println("  - Gray DEBUG, Blue INFO, Yellow WARN, Red ERROR")
	println("✓ File: ./logs/" + time.Now().Format("2006-01-02") + ".log (plain text)")
	println("✓ Daily rotation: New file created each day")
	println("✓ Auto-cleanup: Old logs deleted after 28 days")

	// Example 2: Development Logger - Console Only
	println("\n=== Example 2: Development Logger (Console Only) ===")
	devLogger := logpy.Development()
	devLogger.Debug().Str("component", "auth").Msg("Starting authentication process")
	devLogger.Info().Str("user", "alice").Int("user_id", 123).Msg("User authenticated successfully")
	devLogger.Warn().Str("ip", "192.168.1.1").Int("attempts", 5).Msg("Multiple failed login attempts")
	devLogger.Error().Err(errors.New("connection timeout")).Msg("Failed to connect to database")

	println("\n✓ Console output only (stdout)")
	println("✓ DEBUG level enabled (shows all log levels)")
	println("✓ Perfect for development and debugging")
	println("✓ No file created (console only)")

	// Example 3: File Only - Daily Rotation with Custom Prefix
	println("\n=== Example 3: File Only with Custom Prefix ===")
	customDailyConfig := logpy.Config{
		Level:        logpy.InfoLevel,
		Format:       logpy.FormatConsole,
		Output:       logpy.OutputFile,
		OutputPath:   "./logs/myservice.log", // Creates myservice-YYYY-MM-DD.log
		RotationMode: logpy.RotationDaily,
		UseColor:     false, // Plain text in file
		MaxAge:       7,     // Keep logs for 7 days only
		AddCaller:    true,
		ColorConfig:  logpy.DefaultColorConfig(),
		MultiOutput:  false, // File only, no console
	}
	customLogger := logpy.NewWithConfig(customDailyConfig)
	customLogger.Info().
		Str("service", "api-gateway").
		Int("port", 8080).
		Bool("tls_enabled", true).
		Dur("startup_time", 2*time.Second).
		Msg("Server started")

	println("\n✓ File: ./logs/myservice-" + time.Now().Format("2006-01-02") + ".log")
	println("✓ Plain text (no ANSI color codes)")
	println("✓ No console output")
	println("✓ Retention: 7 days")

	// Example 4: Size-Based Rotation (JSON format)
	println("\n=== Example 4: Size-Based Rotation (JSON Format) ===")
	sizeRotateConfig := logpy.Config{
		Level:        logpy.InfoLevel,
		Format:       logpy.FormatJSON,
		Output:       logpy.OutputFile,
		OutputPath:   "./logs/sizerotate.log", // In logs folder, not /tmp
		RotationMode: logpy.RotationSize,      // Size-based, not daily
		MaxSize:      10,                       // Rotate at 10 MB
		MaxBackups:   3,                        // Keep 3 old files
		MaxAge:       30,                       // Delete after 30 days
		Compress:     true,                     // Compress rotated files (.gz)
		UseColor:     false,                    // No colors for JSON
		MultiOutput:  false,
	}
	sizeLogger := logpy.NewWithConfig(sizeRotateConfig)
	sizeLogger.Info().
		Str("mode", "size-based").
		Str("format", "json").
		Int("max_size_mb", 10).
		Msg("Using lumberjack rotation")

	println("\n✓ JSON format for structured logging")
	println("✓ File: ./logs/sizerotate.log (not /tmp)")
	println("✓ Rotates when file reaches 10MB")
	println("✓ Keeps 3 backup files")
	println("✓ Compresses old files with gzip")

	// Example 5: Child Logger with Persistent Context (Separated Fields)
	println("\n=== Example 5: Child Logger with Context Separation ===")
	requestLogger := devLogger.With(
		logpy.String("request_id", "req-12345"),
		logpy.String("method", "GET"),
		logpy.String("path", "/api/users"),
		logpy.String("client_ip", "192.168.1.100"),
	)

	start := time.Now()
	requestLogger.Info().Msg("Request received")

	// Simulate processing
	time.Sleep(10 * time.Millisecond)

	requestLogger.Info().
		Int("status", 200).
		Dur("latency", time.Since(start)).
		Int("response_size_bytes", 1024).
		Msg("Request completed")

	println("\n✓ Context fields separated with | symbol")
	println("✓ Format: [timestamp] LEVEL file:line message field1=val1 field2=val2 | context1=val1 context2=val2")
	println("✓ Context fields: request_id, method, path, client_ip")
	println("✓ Event fields: status, latency, response_size_bytes")

	// Example 6: All Field Types Demonstration
	println("\n=== Example 6: All Field Types ===")
	devLogger.Info().
		Str("string_field", "hello").
		Int("int_field", 42).
		Int64("int64_field", 9876543210).
		Float64("float_field", 3.14159).
		Bool("bool_field", true).
		Time("time_field", time.Now()).
		Dur("duration_field", 2*time.Second).
		Err(errors.New("example error")).
		Any("custom_field", map[string]string{"key": "value"}).
		Msg("Demonstrating all field types")

	println("\n✓ Field types: Str, Int, Int64, Float64, Bool, Time, Dur, Err, Any")
	println("✓ Type-safe API prevents errors")
	println("✓ Use Any() for custom types (uses reflection)")

	// Example 7: Global Logger
	println("\n=== Example 7: Global Logger ===")
	logpy.SetGlobal(devLogger)

	// Use anywhere in your application without passing logger around
	logpy.Log().Info().Str("global", "true").Msg("Using global logger")
	logpy.Log().Warn().Str("note", "convenient but less testable").Msg("Global logger trade-off")

	println("\n✓ Set once, use anywhere: logpy.Log()")
	println("✓ Convenient for quick logging")
	println("✓ Trade-off: Less testable than dependency injection")

	// Example 8: Custom Color Configuration
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

	println("\n✓ Fully customizable ANSI color codes")
	println("✓ Supports 256-color and RGB colors")
	println("✓ Useful for brand-specific logging themes")

	// Example 9: Production Configuration
	println("\n=== Example 9: Production Logger ===")
	prodLogger := logpy.Production()
	prodLogger.Info().
		Str("environment", "production").
		Str("version", "v1.2.3").
		Msg("Application started in production mode")

	println("\n✓ JSON format for log aggregation (ELK, Splunk, etc.)")
	println("✓ File: /var/log/app.log with rotation")
	println("✓ INFO level (no debug logs)")

	// Example 10: Error Handling with Context
	println("\n=== Example 10: Error Handling with Rich Context ===")
	err := performDatabaseOperation()
	if err != nil {
		devLogger.Error().
			Err(err).
			Str("operation", "database_query").
			Str("table", "users").
			Int("retry_count", 3).
			Dur("timeout", 5*time.Second).
			Msg("Database operation failed")
	}

	println("\n✓ Rich error context for debugging")
	println("✓ Err() field automatically captures error message")
	println("✓ Add operation details for better diagnostics")

	// Summary
	println("\n╔════════════════════════════════════════════════════════════════╗")
	println("║                         SUMMARY                                ║")
	println("╚════════════════════════════════════════════════════════════════╝")

	println("\n📋 Key Configuration Behaviors:")
	println("  1. DEFAULT: Console (colors) + File (plain text)")
	println("  2. DEVELOPMENT: Console only (colors, debug level)")
	println("  3. PRODUCTION: File only (JSON, size-based rotation)")

	println("\n📁 File Naming:")
	println("  - Simple: YYYY-MM-DD.log (e.g., " + time.Now().Format("2006-01-02") + ".log)")
	println("  - Prefix: PREFIX-YYYY-MM-DD.log (e.g., myservice-" + time.Now().Format("2006-01-02") + ".log)")

	println("\n🎨 Color Behavior:")
	println("  - MultiOutput=true:  Console has colors, File is plain")
	println("  - MultiOutput=false: File has no colors (plain text)")

	println("\n📦 Log Files Created:")
	println("  ✓ ./logs/" + time.Now().Format("2006-01-02") + ".log (default, plain)")
	println("  ✓ ./logs/myservice-" + time.Now().Format("2006-01-02") + ".log (plain)")
	println("  ✓ ./logs/sizerotate.log (JSON, size-based)")
	println("  ✓ /var/log/app.log (production)")

	println("\n🔍 View Logs:")
	println("  cat ./logs/" + time.Now().Format("2006-01-02") + ".log")
	println("  tail -f ./logs/" + time.Now().Format("2006-01-02") + ".log")
	println("  grep ERROR ./logs/" + time.Now().Format("2006-01-02") + ".log")

	println("\n✨ All examples completed successfully!")
}

// Helper function for error handling example
func performDatabaseOperation() error {
	return errors.New("connection pool exhausted")
}
