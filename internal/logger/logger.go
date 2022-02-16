// logger provides a lightweight, multi-level logging system, and supports logging to any io.Writer.
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
	LevelDebug Level = -1
	LevelInfo  Level = 0
	LevelWarn  Level = 1
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

type Logger struct {
	debugLog *log.Logger
	infoLog  *log.Logger
	warnLog  *log.Logger
	fatalLog *log.Logger
	Level    Level
}

func New(out io.Writer) Logger {
	return Logger{
		debugLog: log.New(out, colorCyan+"DEBUG \t"+colorReset, log.Ldate|log.Ltime),
		infoLog:  log.New(out, colorBlue+"INFO \t"+colorReset, log.Ldate|log.Ltime),
		warnLog:  log.New(out, colorYellow+"WARN \t"+colorReset, log.Ldate|log.Ltime),
		fatalLog: log.New(out, colorRed+"FATAL \t"+colorReset, log.Ldate|log.Ltime),
	}
}

func (l Logger) Debug(msg string) {
	if l.Level <= LevelDebug {
		l.debugLog.Println(msg)
	}
}

func (l Logger) Debugf(format string, v ...interface{}) {
	if l.Level <= LevelDebug {
		l.debugLog.Printf(format, v...)
	}
}

func (l Logger) Info(msg string) {
	if l.Level <= LevelInfo {
		l.infoLog.Println(msg)
	}
}

func (l Logger) Infof(format string, v ...interface{}) {
	if l.Level <= LevelInfo {
		l.infoLog.Printf(format, v...)
	}
}

func (l Logger) Warn(msg string) {
	if l.Level <= LevelWarn {
		l.warnLog.Println(msg)
	}
}

func (l Logger) Warnf(format string, v ...interface{}) {
	if l.Level <= LevelWarn {
		l.warnLog.Printf(format, v...)
	}
}

func (l Logger) Fatal(msg string) {
	l.fatalLog.Println(msg)
	os.Exit(1)
}

func (l Logger) Fatalf(format string, v ...interface{}) {
	l.fatalLog.Printf(format, v...)
	os.Exit(1)
}
