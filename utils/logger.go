package utils

import (
	"fmt"
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
	Logger.SetOutput(infoLogFile)          // Default output is info.log
	Logger.SetFormatter(&EmojiFormatter{}) // Use the custom emoji formatter
	Logger.SetLevel(logrus.DebugLevel)     // Set the default log level

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

// EmojiFormatter is a custom formatter that adds emojis based on log level
type EmojiFormatter struct{}

func (f *EmojiFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var emoji string
	switch entry.Level {
	case logrus.InfoLevel:
		emoji = "✅" // Emoji for info logs
	case logrus.WarnLevel:
		emoji = "⚠️" // Emoji for warning logs
	case logrus.ErrorLevel:
		emoji = "❌" // Emoji for error logs
	case logrus.DebugLevel:
		emoji = "🐛" // Emoji for debug logs
	case logrus.FatalLevel:
		emoji = "💀" // Emoji for fatal logs
	case logrus.PanicLevel:
		emoji = "🔥" // Emoji for panic logs
	default:
		emoji = "ℹ️" // Default emoji
	}

	// Format the log message with the emoji
	logMessage := fmt.Sprintf("%s [%s] %s\n", emoji, entry.Level.String(), entry.Message)
	return []byte(logMessage), nil
}
