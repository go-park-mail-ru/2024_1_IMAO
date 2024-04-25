package usecases

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(outputPaths []string, errorOutputPaths []string, options ...zap.Option) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.DisableStacktrace = false
	config.OutputPaths = outputPaths
	config.ErrorOutputPaths = errorOutputPaths
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "lvl",
		NameKey:        "",
		CallerKey:      "",
		MessageKey:     "",
		StacktraceKey:  "",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	logger, err := config.Build(options...)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	loggerSugar := logger.Sugar()

	return loggerSugar, nil
}
