package logger

import (
	"github.com/Sterks/fReader/config"
	"github.com/sirupsen/logrus"
)

// Logger ...
type Logger struct {
	logger *logrus.Logger
	config *config.Config
}

// NewLogger ...
func NewLogger() *Logger {
	return &Logger{
		logger: logrus.New(),
		config: &config.Config{},
	}
}

// ConfigureLogger ....
func (l *Logger) ConfigureLogger(conf *config.Config) {
	l.logger = logrus.New()
	cc := conf.MainSettings.LogLevel
	c, err := logrus.ParseLevel(cc)
	if err != nil {
		l.logger.Errorln("Ошибка ", err)
	}
	l.logger.SetLevel(c)
	customFormat := new(logrus.TextFormatter)
	customFormat.TimestampFormat = "2006-01-02 15:04:05"
	customFormat.FullTimestamp = true
	customFormat.ForceColors = true
	l.logger.SetFormatter(customFormat)
}

//InfoLog ...
func (l *Logger) InfoLog(args ...interface{}) {
	l.logger.Infoln(args...)
}

//ErrorLog ...
func (l *Logger) ErrorLog(args ...interface{}) {
	l.logger.Errorln(args...)
}

//DebugLog ...
func (l *Logger) DebugLog(args ...interface{}) {
	l.logger.Debugln(args...)
}
