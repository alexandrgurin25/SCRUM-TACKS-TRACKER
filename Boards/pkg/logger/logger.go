package logger

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

const (
	Key       = "logger"
	RequestID = "request_id"
)

type Logger struct {
	l *zap.Logger
}

func NewLogger(ctx context.Context) (context.Context, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, Key, &Logger{logger})
	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(Key).(*Logger); ok && logger != nil {
		return logger
	}
	return NewDefaultLogger()
}

func NewDefaultLogger() *Logger {
	logger, _ := zap.NewDevelopment()
	return &Logger{l: logger}
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.l.Info(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msq string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.l.Error(msq, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.l.Fatal(msg, fields...)
}

func Interceptor(ctx context.Context) grpc.UnaryServerInterceptor {
	return func(logCtx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		log := GetLoggerFromCtx(ctx)
		logCtx = context.WithValue(logCtx, Key, log)

		md, ok := metadata.FromIncomingContext(logCtx)
		if ok {
			guid, good := md[RequestID]
			if good {
				GetLoggerFromCtx(logCtx).Error(ctx, "no request id")
				ctx = context.WithValue(ctx, RequestID, guid)
			}
		}
		GetLoggerFromCtx(logCtx).Info(logCtx, "request", zap.String("method", info.FullMethod),
			zap.Time("request time", time.Now()),
		)
		return handler(logCtx, req)
	}
}
