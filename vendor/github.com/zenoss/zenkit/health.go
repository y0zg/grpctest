package zenkit

import (
	"context"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type HealthServer interface {
	Serving()
	NotServing()
	Shutdown()
}

func RegisterHealthServer(ctx context.Context, log *logrus.Entry) (HealthServer, error) {
	hs := &healthServer{
		server: grpc.NewServer(
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
				grpc_ctxtags.StreamServerInterceptor(),
				grpc_logrus.StreamServerInterceptor(log),
				grpc_recovery.StreamServerInterceptor(),
			)),
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_logrus.UnaryServerInterceptor(log),
				grpc_recovery.UnaryServerInterceptor(),
			)),
		),
		service: health.NewServer(),
	}
	reflection.Register(hs.server)
	hs.NotServing()
	grpc_health_v1.RegisterHealthServer(hs.server, hs.service)

	addr := viper.GetString(GRPCHealthAddrConfig)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "unable to start health server")
	}
	go hs.server.Serve(lis)
	return hs, nil
}

type healthServer struct {
	server  *grpc.Server
	service *health.Server
}

func (hs *healthServer) Serving() {
	hs.service.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
}

func (hs *healthServer) NotServing() {
	hs.service.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
}

func (hs *healthServer) Shutdown() {
	hs.server.GracefulStop()
}
