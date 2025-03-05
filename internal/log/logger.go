package log

import (
	"context"
	"github.com/DavidMovas/SpeakUp-Server/utils/helpers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

var _ io.Closer = (*Logger)(nil)

type Logger struct {
	*zap.Logger
	file io.Closer
}

func NewLogger(local bool, level string) (*Logger, error) {
	var logger Logger

	logFile, err := os.OpenFile("/var/log/logger.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return nil, err
	}

	logger.file = logFile

	var zapCfg zap.Config
	if local {
		zapCfg = zap.NewDevelopmentConfig()
	} else {
		zapCfg = zap.NewProductionConfig()
	}

	switch level {
	case "debug":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "panic":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	default:
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	zapCfg.OutputPaths = []string{"stdout"}
	zapCfg.ErrorOutputPaths = []string{"stdout"}
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapCfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	zapCfg.DisableStacktrace = true

	logger.Logger, _ = zapCfg.Build(zap.WithCaller(true))

	return &logger, nil
}

func FromContext(ctx context.Context) (*Logger, bool) {
	logger := ctx.Value("logger").(*Logger)
	return logger, logger != nil
}

func (l *Logger) Close() error {
	return helpers.WithClosers([]func() error{l.file.Close, l.Logger.Sync}, nil)
}
