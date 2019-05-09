package zenkit

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/spf13/viper"
	"go.opencensus.io/tag"
	"google.golang.org/grpc"
)

var (
	KeyServiceLabel, _ = tag.NewKey("service_label")
)

func mutators() []tag.Mutator {
	return []tag.Mutator{
		tag.Upsert(KeyServiceLabel, viper.GetString(ServiceLabel)),
	}
}

func MetricTagsStreamServerInterceptor() grpc.StreamServerInterceptor {
	muts := mutators()
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapped := grpc_middleware.WrapServerStream(stream)
		ctx, _ := tag.New(stream.Context(), muts...)
		wrapped.WrappedContext = ctx
		return handler(srv, wrapped)
	}
}

func MetricTagsUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	muts := mutators()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, _ = tag.New(ctx, muts...)
		return handler(ctx, req)
	}
}
