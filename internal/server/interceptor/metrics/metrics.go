package metrics

import (
	"context"
	"path"
	"time"

	"github.com/F0urward/proftwist-backend/internal/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type MetricsUnaryServerInterceptor struct {
	metrics metrics.Metrics
}

func NewMetricsUnaryServerInterceptor(m metrics.Metrics) *MetricsUnaryServerInterceptor {
	return &MetricsUnaryServerInterceptor{
		metrics: m,
	}
}

func (i *MetricsUnaryServerInterceptor) MetricsUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		method := path.Base(info.FullMethod)

		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		code := status.Code(err).String()

		i.metrics.IncGRPCRequest(method, code)
		i.metrics.ObserveGRPCDuration(method, duration)

		if err != nil {
			i.metrics.IncGRPCError(method, code)
		}

		return resp, err
	}
}
