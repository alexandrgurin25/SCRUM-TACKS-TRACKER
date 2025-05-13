package logger

import (
	"context"

	"go.uber.org/zap"
)

const (
	Key       = "logger"
	RequestID = "request_id"
)

type Logger struct {
	l *zap.Logger
}

type ctxKey struct{}

func New(ctx context.Context) (context.Context, error) {
	if GetLoggerFromCtx(ctx) != nil {
		return ctx, nil
	}

	logger, err := zap.NewDevelopment(zap.AddCaller())
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, ctxKey{}, &Logger{logger})
	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	log, ok := ctx.Value(ctxKey{}).(*Logger)
	if !ok || log == nil {
		tmpLogger, _ := zap.NewDevelopment(zap.AddCaller())
		return &Logger{tmpLogger}
	}
	return log
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if requestID := ctx.Value(RequestID); requestID != nil {
		fields = append(fields, zap.String(RequestID, requestID.(string)))

	}
	l.l.Info(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msq string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.l.Error(msq, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msq string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.l.Fatal(msq, fields...)
}
