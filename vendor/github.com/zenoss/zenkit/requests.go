package zenkit

import (
	"context"
	"sync/atomic"

	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"golang.org/x/sync/semaphore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	MConcurrentRequests = stats.Int64("grpc/concurrent_requests", "Concurrent requests", stats.UnitDimensionless)

	ConcurrentRequestsView = &view.View{
		Name:        "grpc/concurrent_requests",
		Description: "Number of concurrent requests",
		Measure:     MConcurrentRequests,
		TagKeys:     []tag.Key{ocgrpc.KeyServerMethod},
		Aggregation: view.LastValue(),
	}
)

func requestCounter(max int) func(context.Context, func() error) error {
	var (
		counter int64
		sem     *semaphore.Weighted
	)
	if max > 0 {
		sem = semaphore.NewWeighted(int64(max))
	}
	return func(ctx context.Context, handler func() error) error {
		if max > 0 {
			if !sem.TryAcquire(1) {
				return status.Error(codes.Unavailable, "concurrent request limit reached")
			}
			defer sem.Release(1)
		}
		stats.Record(ctx, MConcurrentRequests.M(atomic.AddInt64(&counter, 1)))
		defer atomic.AddInt64(&counter, -1)
		return handler()
	}
}

func ConcurrentRequestsStreamServerInterceptor(max int) grpc.StreamServerInterceptor {
	wrapper := requestCounter(max)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return wrapper(stream.Context(), func() error {
			return handler(srv, stream)
		})
	}
}

func ConcurrentRequestsUnaryServerInterceptor(max int) grpc.UnaryServerInterceptor {
	wrapper := requestCounter(max)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = wrapper(ctx, func() error {
			resp, err = handler(ctx, req)
			return err
		})
		return
	}
}
