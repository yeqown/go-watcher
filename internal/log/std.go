package log

import (
	"fmt"
	"os"
)

var (
	std = &stdlogger{logLevel: LevelInfo}
)

// Level ...
type Level uint8

const (
	// LevelFatal ... level of FatalLevel
	LevelFatal Level = iota + 1
	// LevelError ... level of ErrorLevel
	LevelError
	// LevelWarning ... level of WarningLevel
	LevelWarning
	// LevelInfo ... level of InfoLevel
	LevelInfo
	// LevelDebug ... level of DebugLevel
	LevelDebug
)

// Fatal func of (std *stdlogger)
func Fatal(args ...interface{}) {
	std.output(LevelFatal, fmt.Sprint(args...))
	os.Exit(1)
}

// Fatalf func of (std *stdlogger)
func Fatalf(format string, v ...interface{}) {
	std.output(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Error func of (std *stdlogger)
func Error(args ...interface{}) {
	std.output(LevelError, fmt.Sprint(args...))
}

// Errorf func of (std *stdlogger)
func Errorf(format string, v ...interface{}) {
	std.output(LevelError, fmt.Sprintf(format, v...))
}

// Warn func of (std *stdlogger)
func Warn(args ...interface{}) {
	std.output(LevelWarning, fmt.Sprint(args...))
}

// Warnf func of (std *stdlogger)
func Warnf(format string, v ...interface{}) {
	std.output(LevelWarning, fmt.Sprintf(format, v...))
}

// Info func of (std *stdlogger)
func Info(args ...interface{}) {
	std.output(LevelInfo, fmt.Sprint(args...))
}

// Infof func of (std *stdlogger)
func Infof(format string, v ...interface{}) {
	std.output(LevelInfo, fmt.Sprintf(format, v...))
}

// Debug func of (std *stdlogger)
func Debug(args ...interface{}) {
	std.output(LevelDebug, fmt.Sprint(args...))
}

// Debugf func of (std *stdlogger)
func Debugf(format string, v ...interface{}) {
	std.output(LevelDebug, fmt.Sprintf(format, v...))
}

// stdlogger output stdout with fmt
type stdlogger struct {
	logLevel Level
}

func (l *stdlogger) output(level Level, s string) {
	if l.logLevel < level {
		return
	}

	formatStr := "[UNKNOWN] %s"
	switch level {
	case LevelFatal:
		formatStr = "\033[35m[FATAL]\033[0m %s\n"
	case LevelError:
		formatStr = "\033[31m[ERROR]\033[0m %s\n"
	case LevelWarning:
		formatStr = "\033[33m[WARN]\033[0m %s\n"
	case LevelInfo:
		formatStr = "\033[32m[INFO]\033[0m %s\n"
	case LevelDebug:
		formatStr = "\033[36m[DEBUG]\033[0m %s\n"
	}
	fmt.Printf(formatStr, s)
}
