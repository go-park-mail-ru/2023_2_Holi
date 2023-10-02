package logfuncs

import (
	"os"

	"github.com/sirupsen/logrus"
)

func LoggerInit() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})

	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	}
	return logger
}

func LogError(logger *logrus.Logger, packageName, functionName string, err error, message string) {
	logger.WithFields(logrus.Fields{
		"package":  packageName,
		"function": functionName,
		"error":    err,
	}).Error(message)
}

func LogFatal(logger *logrus.Logger, packageName, functionName string, err error, message string) {
	logger.WithFields(logrus.Fields{
		"package":  packageName,
		"function": functionName,
		"error":    err,
	}).Fatal(message)
}
