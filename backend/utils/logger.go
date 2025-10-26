package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is a global logger instance
var Logger *logrus.Logger

// InitLogger initializes the global logger
func InitLogger() {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	Logger.SetLevel(logrus.DebugLevel)
}

// CreateLogger creates a logger with a specific component name
func CreateLogger(component string) *logrus.Entry {
	if Logger == nil {
		InitLogger()
	}
	return Logger.WithField("component", component)
}
