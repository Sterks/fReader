package config

import (
	"github.com/sirupsen/logrus"
)

// Logger ...
type Logger struct {
	logger *logrus.Logger
	config *Config
}

// New ...
func New(config *Config) *Logger {
	return &Logger{
		logger: logrus.New(),
		config: config,
	}
}

// ConfigureLogger ....
func (l *Logger) ConfigureLogger() {
	ll, err := logrus.ParseLevel("debug")
	if err != nil {
		l.logger.Println(err)
	}
	l.logger.SetLevel(ll)
	customFormat := new(logrus.TextFormatter)
	customFormat.TimestampFormat = "2006-01-02 15:04:05"
	customFormat.FullTimestamp = true
	customFormat.ForceColors = true
	l.logger.SetFormatter(customFormat)
}
