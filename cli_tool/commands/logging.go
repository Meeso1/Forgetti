package commands

import (
	"fmt"
)

type LogLevel int

const (
	LogLevelVerbose LogLevel = iota
	LogLevelInfo
	LogLevelError
)

// TODO: log to file too
type Logger interface {
	Verbose(message string, args ...any)
	Info(message string, args ...any)
	Error(message string, args ...any)
}

type ConsoleLogger struct {
	LogLevel LogLevel
}

func MakeLogger(logLevel LogLevel) Logger {
	return &ConsoleLogger{
		LogLevel: logLevel,
	}
}

func (l *ConsoleLogger) Verbose(message string, args ...any) {
	if l.LogLevel <= LogLevelVerbose {
		fmt.Println(fmt.Sprintf(message, args...))
	}
}

func (l *ConsoleLogger) Info(message string, args ...any) {
	if l.LogLevel <= LogLevelInfo {
		fmt.Println(fmt.Sprintf(message, args...))
	}
}

func (l *ConsoleLogger) Error(message string, args ...any) {
	if l.LogLevel <= LogLevelError {
		fmt.Println(fmt.Sprintf(message, args...))
	}
}
