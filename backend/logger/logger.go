package logger

import (
	"os"
	
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	
	// Set log level from environment
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
	
	// Set formatter
	if os.Getenv("LOG_FORMAT") == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	}
	
	// Output to stdout
	log.SetOutput(os.Stdout)
}

// Get returns the logger instance
func Get() *logrus.Logger {
	return log
}

// WithField creates an entry with a single field
func WithField(key string, value interface{}) *logrus.Entry {
	return log.WithField(key, value)
}

// WithFields creates an entry with multiple fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return log.WithFields(fields)
}

// Debug logs at debug level
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Info logs at info level
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warn logs at warn level
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Error logs at error level
func Error(args ...interface{}) {
	log.Error(args...)
}

// Fatal logs at fatal level and exits
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Debugf logs formatted at debug level
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Infof logs formatted at info level
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warnf logs formatted at warn level
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Errorf logs formatted at error level
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatalf logs formatted at fatal level and exits
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
