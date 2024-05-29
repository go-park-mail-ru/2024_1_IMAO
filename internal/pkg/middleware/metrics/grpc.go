package metrics

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/metrics"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	"google.golang.org/grpc"
	"strconv"
	"time"
)

type Interceptor struct {
	metrics metrics.GRPCMetrics
}

func CreateMetricsInterceptor(metrics metrics.GRPCMetrics) *Interceptor {
	return &Interceptor{
		metrics: metrics,
	}
}

func getCode(err string) int {
	switch err {
	case "session does not exist":
		return responses.StatusUnauthorized
	case "no such cookie in userStorage":
		return responses.StatusUnauthorized
	case "user not authorized":
		return responses.StatusUnauthorized
	}

	return responses.StatusBadRequest
}

func (interceptor *Interceptor) ServeMetricsInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	h, err := handler(ctx, req)
	end := time.Since(start)
	code := responses.StatusOk

	if err != nil {
		code = getCode(err.Error())
	}

	codeStr := strconv.Itoa(code)
	interceptor.metrics.IncreaseTotal(codeStr, info.FullMethod)
	interceptor.metrics.AddDuration(codeStr, info.FullMethod, end)

	return h, err
}
