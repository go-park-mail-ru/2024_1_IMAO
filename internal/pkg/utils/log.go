package utils

import (
	"context"
	"runtime"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/config"

	newlogger "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/usecases"
)

const (
	outputLogPath      = "stdout logs.json"
	errorOutputLogPath = "stderr err_logs.json"
)

func GFN() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	values := strings.Split(frame.Function, "/")

	return values[len(values)-1]
}

func GetRequestId(ctx context.Context) string {
	requestID, _ := ctx.Value(config.RequestUUIDContextKey).(uuid.UUID)
	return requestID.String()
}

func LogHandlerInfo(logger *zap.SugaredLogger, statusCode int, msg string) {
	logger = logger.With(zap.String("status", strconv.Itoa(statusCode)))
	logger.Info(msg)
}

func LogHandlerError(logger *zap.SugaredLogger, statusCode int, msg string) {
	logger = logger.With(zap.String("status", strconv.Itoa(statusCode)))
	logger.Error(msg)
}

func GetLoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(config.LoggerContextKey).(*zap.SugaredLogger); ok {
		return logger
	}

	logger, _ := newlogger.NewLogger(strings.Split(outputLogPath, " "), strings.Split(errorOutputLogPath, " "))
	logger.Error("Couldnt get logger from context")

	return logger
}
