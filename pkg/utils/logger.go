package utils

import (
	"fmt"
	"time"
)

type LogLevel int

const (
	ErrorLevel LogLevel = iota
	InfoLevel
	DebugLevel
)

type Logger struct {
	verbose bool
}

func NewLogger(verbose bool) *Logger {
	return &Logger{
		verbose: verbose,
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.verbose {
		l.log(DebugLevel, format, args...)
	}
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(InfoLevel, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ErrorLevel, format, args...)
}

func (l *Logger) CommandOutput(command string, output string) {
	if output == "" {
		return
	}

	fmt.Printf("\n--- Command Output: %s ---\n", command)
	fmt.Println(output)
	fmt.Println("--- End of Output ---\n")
}

func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	timeStr := time.Now().Format("15:04:05")
	var levelStr string

	switch level {
	case DebugLevel:
		levelStr = "\033[36mDEBUG\033[0m" // Cyan
	case InfoLevel:
		levelStr = "\033[32mINFO\033[0m" // Green
	case ErrorLevel:
		levelStr = "\033[31mERROR\033[0m" // Red
	}

	message := fmt.Sprintf(format, args...)
	fmt.Printf("[%s] %s: %s\n", timeStr, levelStr, message)
}
