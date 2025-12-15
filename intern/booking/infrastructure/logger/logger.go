// internal/infrastructure/logger/slog_logger.go
package logger

import (
	"final-project/pkg/booking/constants"
	"log/slog"
	"os"
)

type SlogLogger struct {
	log *slog.Logger
}

func New(level slog.Level) *SlogLogger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return &SlogLogger{
		log: slog.New(handler),
	}
}

func (l *SlogLogger) Debug(msg string, fields ...any) {
	l.log.Debug(msg, fields...)
}

func (l *SlogLogger) Info(msg string, fields ...any) {
	l.log.Info(msg, fields...)
}

func (l *SlogLogger) Warn(msg string, fields ...any) {
	l.log.Warn(msg, fields...)
}

func (l *SlogLogger) Error(msg string, err error, fields ...any) {
	fields = append(fields, constants.KeyError, err)
	l.log.Error(msg, fields...)
}
