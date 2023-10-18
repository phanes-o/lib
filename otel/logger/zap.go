package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLog struct {
	logger *otelzap.Logger
	opts *Options
}

func (z *zapLog) Sync() {
	z.logger.Sync()
}

func (z *zapLog) WithFields(fields ...zapcore.Field) Logger {
	newLogger := z.logger.With(fields...)
	newZap := &zapLog{
		logger: otelzap.New(newLogger),
		opts: z.opts,
	}
	return newZap
}

func (z *zapLog) Ctx(ctx context.Context) Logger {
	newCtx := z.logger.Ctx(ctx)
	newZap := &zapLog{
		logger: newCtx.Logger(),
		opts: z.opts,
	}
	return newZap
}

func (z *zapLog) Debug(msg string, fields ...zapcore.Field) {
	z.logger.Debug(z.withPrefix(msg), fields...)
}

func (z *zapLog) DebugCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.DebugContext(ctx, z.withPrefix(msg), fields...)
}

func (z *zapLog) Info(msg string, fields ...zapcore.Field) {
	z.logger.Info(z.withPrefix(msg), fields...)
}

func (z *zapLog) InfoCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.InfoContext(ctx, z.withPrefix(msg), fields...)
}

func (z *zapLog) Warn(msg string, fields ...zapcore.Field) {
	z.logger.Warn(z.withPrefix(msg), fields...)
}

func (z *zapLog) WarnCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.WarnContext(ctx, z.withPrefix(msg), fields...)
}

func (z *zapLog) Error(msg string, fields ...zapcore.Field) {
	z.logger.Error(z.withPrefix(msg), fields...)
}

func (z *zapLog) ErrorCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.ErrorContext(ctx, z.withPrefix(msg), fields...)
}

func (z *zapLog) Fatal(msg string, fields ...zapcore.Field) {
	z.logger.Fatal(z.withPrefix(msg), fields...)
}

func (z *zapLog) FatalCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.FatalContext(ctx, z.withPrefix(msg), fields...)
}

func (z *zapLog) PanicCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	z.logger.PanicContext(ctx, z.withPrefix(msg), fields...)
}

func (z *zapLog) withPrefix(msg string) string {
	return z.opts.prefix + " " + msg
}

func NewZapLog(opts ...Option) Logger {
	options := &Options{
		skip: 2,
		level: -1,
		interval: 10,
		bufferSize: 4086,
		prefix: "",
		stdout: true,
		writers: []io.Writer{os.Stderr},
	}

	for _, o := range opts {
		o.Apply(options)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder

	syncer := newBufferedWriteSyncer(
		options.bufferSize,
		time.Duration(options.interval)*time.Second,
		io.MultiWriter(options.writers...),
	)

	outputCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), syncer, zapcore.Level(options.level))
	core := zapcore.NewTee(outputCore)

	if options.stdout {
		consoleCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), os.Stdout, zapcore.Level(options.level))
		core = zapcore.NewTee(consoleCore, outputCore)
	}

	traceLogger := otelzap.New(
		zap.New(core, zap.AddCaller(), zap.AddCallerSkip(options.skip)),
		otelzap.WithCaller(true),
		otelzap.WithTraceIDField(true),
	)

	logger := &zapLog{
		logger: traceLogger,
		opts: options,
	}
	return logger
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
