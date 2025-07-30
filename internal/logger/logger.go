package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Logger defines the interface for the logging system.
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
	Warning(format string, v ...interface{})
}

type defaultLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	warnLogger  *log.Logger
}

var (
	loggerInstance Logger
	once           sync.Once
)

// NewLogger creates and returns a new instance of defaultLogger.
// It configures the loggers to write to both the console and a file.
func NewLogger() Logger {
	once.Do(func() {
		logDir := "logs"
		logFileName := "app.log"
		logFilePath := filepath.Join(logDir, logFileName)

		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			if err = os.Mkdir(logDir, 0755); err != nil {
				log.Printf("ERROR: Could not create logs directory '%s': %v. Logs will be console-only.", logDir, err)
				loggerInstance = &defaultLogger{
					infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
					errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
					fatalLogger: log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
					warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
				}
				return
			}
		}

		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("ERROR: Could not open log file '%s': %v. Logs will be console-only.", logFilePath, err)
			loggerInstance = &defaultLogger{
				infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
				errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
				fatalLogger: log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
				warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
			}
			return
		}

		infoWarnMultiWriter := io.MultiWriter(os.Stdout, file)
		errorFatalMultiWriter := io.MultiWriter(os.Stderr, file)

		loggerInstance = &defaultLogger{
			infoLogger:  log.New(infoWarnMultiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
			errorLogger: log.New(errorFatalMultiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
			fatalLogger: log.New(errorFatalMultiWriter, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
			warnLogger:  log.New(infoWarnMultiWriter, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		}
	})
	return loggerInstance
}

// Info implements the Info method for defaultLogger.
func (l *defaultLogger) Info(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

// Error implements the Error method for defaultLogger.
func (l *defaultLogger) Error(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

// Fatal implements the Fatal method for defaultLogger.
func (l *defaultLogger) Fatal(format string, v ...interface{}) {
	l.fatalLogger.Fatalf(format, v...) // Exits the application after logging
}

// Warning implements the Warning method for defaultLogger.
func (l *defaultLogger) Warning(format string, v ...interface{}) {
	l.warnLogger.Printf(format, v...)
}
