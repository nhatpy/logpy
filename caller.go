package logpy

import (
	"path/filepath"
	"runtime"
)

// CallerInfo contains information about where a log was called from
type CallerInfo struct {
	File     string
	Line     int
	Function string
}

// getCaller retrieves caller information from the call stack
// skip is the number of stack frames to skip (typically 2-4 depending on call depth)
func getCaller(skip int) CallerInfo {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return CallerInfo{
			File:     "unknown",
			Line:     0,
			Function: "unknown",
		}
	}

	// Get function name
	fn := runtime.FuncForPC(pc)
	funcName := "unknown"
	if fn != nil {
		funcName = fn.Name()
	}

	return CallerInfo{
		File:     filepath.Base(file), // Only keep filename, not full path
		Line:     line,
		Function: funcName,
	}
}
