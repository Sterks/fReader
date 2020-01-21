package logger

import (
	"github.com/Sterks/FReader/config"
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
	logr := logrus.New()
	l.logger = logr
	cc := conf.MainSettings.LogLevel
	c, err := logrus.ParseLevel(cc)
	if err != nil {
		l.logger.Errorf("Ошибка %v", err)
	}
	l.logger.SetLevel(c)
	customFormat := new(logrus.TextFormatter)
	customFormat.TimestampFormat = "2006-01-02 15:04:05"
	customFormat.FullTimestamp = true
	customFormat.ForceColors = true
	l.logger.SetFormatter(customFormat)
}

//InfoLog ...
func (l *Logger) InfoLog(mes string) {
	l.logger.Info(mes)
}

//ErrorLog ...
func (l *Logger) ErrorLog(mes string, err error) {
	l.logger.Errorf(mes, err)
}

//DebugLog ...
func (l *Logger) DebugLog(mes string) {
	l.logger.Debug(mes)
}
