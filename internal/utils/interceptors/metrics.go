package interceptors

import (
	"context"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
)

func NewChainUnaryRequestsCounterInterceptor(rpc metric.Int64Counter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		rpc.Add(ctx, 1)

		return handler(ctx, req)
	}
}
