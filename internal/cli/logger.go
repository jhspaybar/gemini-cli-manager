package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var debugLogger *log.Logger

func init() {
	// Create logs directory if it doesn't exist
	logsDir := filepath.Join(".", "logs")
	os.MkdirAll(logsDir, 0755)

	// Create log file with timestamp
	logFile := filepath.Join(logsDir, fmt.Sprintf("debug-%s.log", time.Now().Format("2006-01-02-15-04-05")))

	file, err := os.Create(logFile)
	if err != nil {
		// Fallback to stderr if we can't create log file
		debugLogger = log.New(os.Stderr, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
		return
	}

	debugLogger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger.Printf("=== Gemini CLI Manager Debug Log Started ===")
}

// LogDebug logs a debug message
func LogDebug(format string, v ...interface{}) {
	if debugLogger != nil {
		debugLogger.Printf(format, v...)
	}
}

// LogMessage logs a Bubble Tea message
func LogMessage(source string, msg interface{}) {
	if debugLogger != nil {
		debugLogger.Printf("[%s] Message: %T %+v", source, msg, msg)
	}
}
