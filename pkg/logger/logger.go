package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	With(fields ...Field) Logger
}

type zapLogger struct {
	log *zap.Logger
}

type Field struct {
	Key   string
	Value any
}

func New(lvl string) (Logger, error) {
	var cfg zap.Config

	if lvl == "debug" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	log, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &zapLogger{
		log: log,
	}, nil
}

func (z *zapLogger) Info(msg string, fields ...Field) {
	z.log.Info(msg, toZap(fields)...)
}

func (z *zapLogger) Error(msg string, fields ...Field) {
	z.log.Error(msg, toZap(fields)...)
}

func (z *zapLogger) Debug(msg string, fields ...Field) {
	z.log.Debug(msg, toZap(fields)...)
}

func (z *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{
		log: z.log.With(toZap(fields)...),
	}
}

func toZap(fields []Field) []zap.Field {
	temple := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		temple = append(temple, zap.Any(f.Key, f.Value))
	}

	return temple
}
