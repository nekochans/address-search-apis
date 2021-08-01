package infrastructure

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zapLogger       *zap.Logger
	lambdaRequestId string
	httpRequestId   string
}

func CreateLogger(lambdaReqId string, httpReqId string) *Logger {
	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.DebugLevel)

	zapConfig := zap.Config{
		Level:    level,
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "name",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stackTrace",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	zapLogger, _ := zapConfig.Build(zap.AddCallerSkip(1))

	logger := &Logger{
		zapLogger:       zapLogger,
		lambdaRequestId: lambdaReqId,
		httpRequestId:   httpReqId,
	}

	return logger
}

func (l *Logger) Info(msg string) {
	l.zapLogger.Info(
		msg,
		zap.String("lambdaRequestId", l.lambdaRequestId),
		zap.String("httpRequestId", l.httpRequestId),
	)
}

func (l *Logger) Error(msg string) {
	l.zapLogger.Error(
		msg,
		zap.String("lambdaRequestId", l.lambdaRequestId),
		zap.String("httpRequestId", l.httpRequestId),
	)
}
