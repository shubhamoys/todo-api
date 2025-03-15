package utils

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func InitLogger() {
	// Create the logs directory if it doesn't exist
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755)
		if err != nil {
			logrus.Fatalf("Failed to create logs directory: %v", err)
		}
	}

	// Create separate log files for info and error logs
	infoLogFile, err := os.OpenFile(filepath.Join(logDir, "info.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open info log file: %v", err)
	}

	errorLogFile, err := os.OpenFile(filepath.Join(logDir, "error.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open error log file: %v", err)
	}

	// Set the output for different log levels
	Logger.SetOutput(infoLogFile) // Default output is info.log
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetLevel(logrus.DebugLevel) // Set the default log level

	// Add hooks to write error logs to error.log
	Logger.AddHook(&LogFileHook{
		Writer:    errorLogFile,
		LogLevels: []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
	})
}

// LogFileHook is a custom hook to write logs to specific files
type LogFileHook struct {
	Writer    *os.File
	LogLevels []logrus.Level
}

func (hook *LogFileHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

func (hook *LogFileHook) Levels() []logrus.Level {
	return hook.LogLevels
}
