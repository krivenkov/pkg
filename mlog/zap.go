package mlog

import (
	"context"

	"github.com/krivenkov/pkg/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level  string `short:"l" long:"level" env:"LEVEL" description:"logging level" default:"INFO" json:"level" yaml:"level"`
	Encode string `long:"encode" env:"ENCODE" description:"log as json or not" default:"json" json:"encode" yaml:"encode"`
	Caller string `long:"caller" env:"CALLER" description:"tell mlog how to encode caller path (short or full)" default:"" json:"caller" yaml:"caller"`
}

const (
	collerFull  = "full"
	collerShort = "short"
)

func New(c *Config) (*zap.Logger, error) {
	var (
		err    error
		logger *zap.Logger
	)
	config := zap.NewProductionConfig()
	config.Encoding = c.Encode
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	config.EncoderConfig.TimeKey = "ts"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.FunctionKey = zapcore.OmitKey
	config.EncoderConfig.MessageKey = "msg"
	config.EncoderConfig.StacktraceKey = "stacktrace"
	config.EncoderConfig.LineEnding = zapcore.DefaultLineEnding
	config.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	config.OutputPaths = []string{"stdout"}

	if c.Encode == "json" && c.Caller == "" {
		c.Caller = collerShort
	}

	switch c.Caller {
	case collerFull:
		config.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	case collerShort:
		config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	default:
		config.EncoderConfig.EncodeCaller = nil
	}

	if errLocal := config.Level.UnmarshalText([]byte(c.Level)); errLocal != nil && c.Level != "" {
		return nil, errLocal
	}
	logger, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func NewC(c Config, info global.Info) (*zap.Logger, error) {
	l, err := New(&c)
	if err != nil {
		return nil, err
	}
	return l.With(zap.String("app", info.AppName)), nil
}

type ctxKey struct{} // or exported to use outside the package

func CtxWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	if ctxLogger, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return ctxLogger
	}
	return zap.NewNop()
}
