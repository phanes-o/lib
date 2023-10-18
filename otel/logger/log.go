package logger

import (
	"context"

	"go.uber.org/zap/zapcore"
)

type Logger interface {
	WithFields(fields ...zapcore.Field) Logger
	Ctx(ctx context.Context) Logger
	Debug(msg string, fields ...zapcore.Field)
	DebugCtx(ctx context.Context, msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	InfoCtx(ctx context.Context, msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	WarnCtx(ctx context.Context, msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	ErrorCtx(ctx context.Context, msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	FatalCtx(ctx context.Context, msg string, fields ...zapcore.Field)
	PanicCtx(ctx context.Context, msg string, fields ...zapcore.Field)
	Sync()
}
