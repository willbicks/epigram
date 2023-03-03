// Package logger provides a lightweight, multi-level logging system, and supports logging to any io.Writer.
package logger

import (
	"io"
	"log"
	"os"
)

// Level represents a level of logging severity / detail. When an desired log level is selected, all messages at or
// above that level are printed, while all levels below are suppressed.
type Level int8

// Log Levels
const (
	// LevelDebug is the lowest level of logging, and is intended for debugging purposes.
	LevelDebug Level = -1
	// LevelInfo is intended for general information including records of operations performed.
	LevelInfo Level = 0
	// LevelWarn is intended for messages that indicate some unexpected non-fatal error has ocured.
	LevelWarn Level = 1
	// LevelFatal is the highest level of logging, and is intended for messages that indicate a fatal error.
	LevelFatal Level = 2
)

// ANSI color escape codes
const (
	colorReset = "\033[0m"

	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

// Logger provides methods to log messages by Level
type Logger struct {
	debugLog *log.Logger
	infoLog  *log.Logger
	warnLog  *log.Logger
	fatalLog *log.Logger
	Level    Level
}

// New returns a new Logger that writes to the provided io.Writer. If timestamp is true, the logger will prefix each
// message with a timestamp.
//
// If the NO_COLOR environment variable is set, the logger will not use ANSI color escape codes.
	var flags int
	if timestamp {
		flags = log.Ldate | log.Ltime
	}

	if os.Getenv("NO_COLOR") != "" {
		return Logger{
			debugLog: log.New(out, "DEBUG \t", flags),
			infoLog:  log.New(out, "INFO \t", flags),
			warnLog:  log.New(out, "WARN \t", flags),
			fatalLog: log.New(out, "FATAL \t", flags),
		}

	}

	return Logger{
		debugLog: log.New(out, colorCyan+"DEBUG \t"+colorReset, flags),
		infoLog:  log.New(out, colorBlue+"INFO \t"+colorReset, flags),
		warnLog:  log.New(out, colorYellow+"WARN \t"+colorReset, flags),
		fatalLog: log.New(out, colorRed+"FATAL \t"+colorReset, flags),
	}
}

// Debug logs a message at the Debug level
func (l Logger) Debug(msg string) {
	if l.Level <= LevelDebug {
		l.debugLog.Println(msg)
	}
}

// Debugf logs a formatted message at the Debug level
func (l Logger) Debugf(format string, v ...interface{}) {
	if l.Level <= LevelDebug {
		l.debugLog.Printf(format, v...)
	}
}

// Info logs a message at the Info level
func (l Logger) Info(msg string) {
	if l.Level <= LevelInfo {
		l.infoLog.Println(msg)
	}
}

// Infof logs a formatted message at the Info level
func (l Logger) Infof(format string, v ...interface{}) {
	if l.Level <= LevelInfo {
		l.infoLog.Printf(format, v...)
	}
}

// Warn logs a message at the Warn level
func (l Logger) Warn(msg string) {
	if l.Level <= LevelWarn {
		l.warnLog.Println(msg)
	}
}

// Warnf logs a formatted message at the Warn level
func (l Logger) Warnf(format string, v ...interface{}) {
	if l.Level <= LevelWarn {
		l.warnLog.Printf(format, v...)
	}
}

// Fatal logs a message at the Fatal level and exits the program with status code 1
func (l Logger) Fatal(msg string) {
	l.fatalLog.Println(msg)
	os.Exit(1)
}

// Fatalf logs a formatted message at the Fatal level and exits the program with status code 1
func (l Logger) Fatalf(format string, v ...interface{}) {
	l.fatalLog.Printf(format, v...)
	os.Exit(1)
}
