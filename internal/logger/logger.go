package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Logger defines the interface for the logging system.
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
	Warning(format string, v ...interface{})
}

type LogEntry struct {
	Level     string `json:"level"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	File      string `json:"file,omitempty"`
	Line      int    `json:"line,omitempty"`
}

type defaultLogger struct {
	infoWriter  io.Writer
	errorWriter io.Writer
	fatalWriter io.Writer
	warnWriter  io.Writer
	stdLogger   *log.Logger
}

var (
	loggerInstance Logger
	once           sync.Once
)

// NewLogger creates and returns a new instance of defaultLogger.
func NewLogger() Logger {
	once.Do(func() {
		logDir := "logs"
		logFileName := "app.log"
		logFilePath := filepath.Join(logDir, logFileName)

		internalStdLogger := log.New(os.Stderr, "LOGGER_ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			if err = os.Mkdir(logDir, 0755); err != nil {
				internalStdLogger.Printf("Could not create logs directory '%s': %v. Logs will be console-only.", logDir, err)
				loggerInstance = createConsoleOnlyLogger(internalStdLogger)
				return
			}
		}

		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			internalStdLogger.Printf("Could not open log file '%s': %v. Logs will be console-only.", logFilePath, err)
			loggerInstance = createConsoleOnlyLogger(internalStdLogger)
			return
		}

		infoWarnMultiWriter := io.MultiWriter(os.Stdout, file)
		errorFatalMultiWriter := io.MultiWriter(os.Stderr, file)

		loggerInstance = &defaultLogger{
			infoWriter:  infoWarnMultiWriter,
			errorWriter: errorFatalMultiWriter,
			fatalWriter: errorFatalMultiWriter,
			warnWriter:  infoWarnMultiWriter,
			stdLogger:   internalStdLogger,
		}
	})
	return loggerInstance
}

// createConsoleOnlyLogger provides a fallback if file logging fails.
func createConsoleOnlyLogger(internalStdLogger *log.Logger) *defaultLogger {
	return &defaultLogger{
		infoWriter:  os.Stdout,
		errorWriter: os.Stderr,
		fatalWriter: os.Stderr,
		warnWriter:  os.Stdout,
		stdLogger:   internalStdLogger,
	}
}

// writeLogEntry formats the log entry as JSON and writes it to the given writer.
func (l *defaultLogger) writeLogEntry(level, format string, writer io.Writer, v ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	var callerFile string
	var callerLine int
	if ok {
		callerFile = filepath.Base(file)
		callerLine = line
	}

	entry := LogEntry{
		Level:     level,
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   fmt.Sprintf(format, v...),
		File:      callerFile,
		Line:      callerLine,
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		l.stdLogger.Printf("Failed to marshal log entry to JSON: %v. Original message: %s", err, fmt.Sprintf(format, v...))
		fmt.Fprintf(writer, "[%s] [%s] %s:%d %s\n", level, time.Now().Format(time.RFC3339), callerFile, callerLine, fmt.Sprintf(format, v...))
		return
	}

	_, err = writer.Write(append(jsonBytes, '\n'))
	if err != nil {
		l.stdLogger.Printf("Failed to write log entry: %v. Original message: %s", err, fmt.Sprintf(format, v...))
	}
}

// Info implements the Info method for defaultLogger.
func (l *defaultLogger) Info(format string, v ...interface{}) {
	l.writeLogEntry("INFO", format, l.infoWriter, v...)
}

// Error implements the Error method for defaultLogger.
func (l *defaultLogger) Error(format string, v ...interface{}) {
	l.writeLogEntry("ERROR", format, l.errorWriter, v...)
}

// Fatal implements the Fatal method for defaultLogger.
func (l *defaultLogger) Fatal(format string, v ...interface{}) {
	l.writeLogEntry("FATAL", format, l.fatalWriter, v...)
	os.Exit(1) // Fatal logs should exit the application
}

// Warning implements the Warning method for defaultLogger.
func (l *defaultLogger) Warning(format string, v ...interface{}) {
	l.writeLogEntry("WARN", format, l.warnWriter, v...)
}
