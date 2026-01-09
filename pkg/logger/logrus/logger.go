package logrus

import (
	"context"

	"github.com/F0urward/proftwist-backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	entry *logrus.Entry
}

func NewLogrusLogger() *LogrusLogger {
	log := logrus.New()

	return &LogrusLogger{
		entry: logrus.NewEntry(log),
	}
}

func (l *LogrusLogger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *LogrusLogger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *LogrusLogger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *LogrusLogger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *LogrusLogger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *LogrusLogger) Panic(args ...interface{}) {
	l.entry.Panic(args...)
}

func (l *LogrusLogger) Panicf(format string, args ...interface{}) {
	l.entry.Panicf(format, args...)
}

func (l *LogrusLogger) WithField(key string, value interface{}) logger.Logger {
	return &LogrusLogger{
		entry: l.entry.WithField(key, value),
	}
}

func (l *LogrusLogger) WithFields(fields map[string]interface{}) logger.Logger {
	return &LogrusLogger{
		entry: l.entry.WithFields(fields),
	}
}

func (l *LogrusLogger) WithError(err error) logger.Logger {
	return &LogrusLogger{
		entry: l.entry.WithError(err),
	}
}

func (l *LogrusLogger) WithContext(ctx context.Context) logger.Logger {
	return &LogrusLogger{
		entry: l.entry.WithContext(ctx),
	}
}
