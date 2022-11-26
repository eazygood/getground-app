package log

import (
	"github.com/eazygood/getground-app/internal/config"
	logger "github.com/sirupsen/logrus"
)

// Init initializes log for the service from the provided config
func Init(cfg config.Log) {
	setLogFormatter(cfg.Formatter)
	setLogLevel(cfg.Level)
}

// setLogFormatter sets the log formatter (Defaults to JSON)
func setLogFormatter(formatter string) {
	switch formatter {
	case "", "JSON":
		logger.SetFormatter(&logger.JSONFormatter{})
	case "TEXT":
		logger.SetFormatter(&logger.TextFormatter{})
	default:
		logger.SetFormatter(&logger.JSONFormatter{})
		logger.Errorf("Invalid logging formatter %s: Using JSON", formatter)
	}
}

// setLogLevel sets tje log level for the application (Defaults to Info)
func setLogLevel(level string) {
	if len(level) == 0 {
		level = "INFO"
	}

	l, err := logger.ParseLevel(level)
	if err != nil {
		logger.WithError(err).Error("Invalid logging level: Using INFO")
		logger.SetLevel(logger.InfoLevel)
		return
	}

	logger.Infof("Setting the log level to %s", level)
	logger.SetLevel(l)
}
