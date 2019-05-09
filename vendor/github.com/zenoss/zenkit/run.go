package zenkit

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opencensus.io/tag"
	"google.golang.org/grpc"
)

type ServiceRegistrationFunc func(*grpc.Server) error

func RunGRPCServerWithTLS(ctx context.Context, name string, f ServiceRegistrationFunc) error {
	return runGRPCServer(ctx, name, f, true, false)
}

func RunGRPCServerWithHealth(ctx context.Context, name string, f ServiceRegistrationFunc) error {
	return runGRPCServer(ctx, name, f, false, true)
}

func RunGRPCServer(ctx context.Context, name string, f ServiceRegistrationFunc) error {
	return runGRPCServer(ctx, name, f, false, false)
}

func WithTrapSIGINT(ctx context.Context, log *logrus.Entry) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case <-c:
			log.Info("SIGINT received")
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx
}

func runGRPCServer(ctx context.Context, name string, f ServiceRegistrationFunc, useTLS bool, startHealth bool) error {
	// If useTLS is true, grpc will listen using tls cert.
	// If grpc based startHealth is true and useTLS is true, the grpc health will not use tls.

	InitConfig(name)

	// Initialize tags from config
	ctx, _ = tag.New(ctx,
		tag.Upsert(KeyServiceLabel, viper.GetString(ServiceLabel)),
	)

	log := Logger(name)
	grpc_logrus.ReplaceGrpcLogger(log)

	// Start listening for SIGINT
	ctx = WithTrapSIGINT(ctx, log)

	if startHealth {
		// Set up the health server
		health, err := RegisterHealthServer(ctx, log)
		if err != nil {
			return err
		}
		health.Serving()
		defer health.NotServing()
		defer health.Shutdown()
	}

	// Create a GRPC server
	server := NewGRPCServer(ctx, name, useTLS, log)
	if err := f(server); err != nil {
		return errors.Wrap(err, "unable to register service")
	}

	// Create a listener
	addr := viper.GetString(GRPCListenAddrConfig)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "unable to start listener")
	}

	// Let's go!
	go server.Serve(lis)
	log.WithField("address", addr).Info("started server")

	// Wait for a cancel from upstream or a SIGINT
	<-ctx.Done()

	server.GracefulStop()
	log.Info("shut down server gracefully")

	return nil
}
