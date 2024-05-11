package metrics

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/metrics"
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
		return 401
	case "no such cookie in userStorage":
		return 401
	case "user not authorized":
		return 401
	}

	return 400
}

func (interceptor *Interceptor) ServeMetricsInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	h, err := handler(ctx, req)
	end := time.Since(start)
	code := 200
	if err != nil {
		code = getCode(err.Error())
	}

	codeStr := strconv.Itoa(code)
	interceptor.metrics.IncreaseTotal(codeStr, info.FullMethod)
	interceptor.metrics.AddDuration(codeStr, info.FullMethod, end)
	return h, err
}
