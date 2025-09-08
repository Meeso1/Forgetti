package logging

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type LogLevel int

const (
	LogLevelVerbose LogLevel = iota
	LogLevelInfo
	LogLevelError
)

type Logger interface {
	Verbose(message string, args ...any)
	Info(message string, args ...any)
	Error(message string, args ...any)
}

type Config struct {
	LogLevel LogLevel
	LogFile  string // if empty, logs only to console
}

type MultiLogger struct {
	LogLevel LogLevel
	Context  string // context information (filename/struct name)
	writers  []io.Writer
	mutex    sync.Mutex
}

var (
	globalConfig Config
	configMutex  sync.RWMutex
	logFile      *os.File
)

// SetGlobalConfig sets the global logging configuration
func SetGlobalConfig(config Config) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	// Close existing log file if any
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}

	globalConfig = config

	// Open new log file if specified
	if config.LogFile != "" {
		var err error
		logFile, err = os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file %s: %w", config.LogFile, err)
		}
	}

	return nil
}

// GetGlobalConfig returns a copy of the global logging configuration
func GetGlobalConfig() Config {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return globalConfig
}

// MakeLogger creates a new logger with the specified context
func MakeLogger(context string) Logger {
	configMutex.RLock()
	defer configMutex.RUnlock()

	writers := []io.Writer{os.Stdout}
	if logFile != nil {
		writers = append(writers, logFile)
	}

	return &MultiLogger{
		LogLevel: globalConfig.LogLevel,
		Context:  context,
		writers:  writers,
	}
}

// MakeLoggerWithLevel creates a new logger with the specified context and log level override
func MakeLoggerWithLevel(context string, level LogLevel) Logger {
	configMutex.RLock()
	defer configMutex.RUnlock()

	writers := []io.Writer{os.Stdout}
	if logFile != nil {
		writers = append(writers, logFile)
	}

	return &MultiLogger{
		LogLevel: level,
		Context:  context,
		writers:  writers,
	}
}

func (l *MultiLogger) log(level string, message string, args ...any) {
	if len(l.writers) == 0 {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	contextStr := ""
	if l.Context != "" {
		contextStr = fmt.Sprintf("[%s] ", l.Context)
	}

	logMessage := fmt.Sprintf("[%s] %s%s%s\n", timestamp, level, contextStr, fmt.Sprintf(message, args...))

	for _, writer := range l.writers {
		writer.Write([]byte(logMessage))
	}
}

func (l *MultiLogger) Verbose(message string, args ...any) {
	if l.LogLevel <= LogLevelVerbose {
		l.log("VERBOSE ", message, args...)
	}
}

func (l *MultiLogger) Info(message string, args ...any) {
	if l.LogLevel <= LogLevelInfo {
		l.log("INFO ", message, args...)
	}
}

func (l *MultiLogger) Error(message string, args ...any) {
	if l.LogLevel <= LogLevelError {
		l.log("ERROR ", message, args...)
	}
}
