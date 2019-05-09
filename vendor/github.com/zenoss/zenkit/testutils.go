package zenkit

import (
	"math/rand"
	"net"
	"time"

	grpc_testing "github.com/grpc-ecosystem/go-grpc-middleware/testing"
	pb_testproto "github.com/grpc-ecosystem/go-grpc-middleware/testing/testproto"
	"github.com/onsi/ginkgo"
	"google.golang.org/grpc"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type GRPCTestServer struct {
	T           ginkgo.GinkgoTInterface
	TestService pb_testproto.TestServiceServer
	ServerOpts  []grpc.ServerOption
	ClientOpts  []grpc.DialOption
	Listener    net.Listener
	Server      *grpc.Server
	clientConn  *grpc.ClientConn
	Client      pb_testproto.TestServiceClient
}

func (s *GRPCTestServer) Setup() {
	var err error
	s.Listener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		s.T.Error(err, "unable to assign test server listener port")
	}
	s.Server = grpc.NewServer(s.ServerOpts...)
	if s.TestService == nil {
		s.TestService = &grpc_testing.TestPingService{}
	}
	pb_testproto.RegisterTestServiceServer(s.Server, s.TestService)
	go func() {
		s.Server.Serve(s.Listener)
	}()
	s.Client = s.NewClient(s.ClientOpts...)
}

func (s *GRPCTestServer) NewClient(opts ...grpc.DialOption) pb_testproto.TestServiceClient {
	newDialOpts := append(opts, grpc.WithBlock(), grpc.WithTimeout(2*time.Second), grpc.WithInsecure())
	clientConn, err := grpc.Dial(s.ServerAddr(), newDialOpts...)
	if err != nil {
		s.T.Error(err, "test client unable to dial test server")
	}
	s.clientConn = clientConn
	return pb_testproto.NewTestServiceClient(clientConn)
}

func (s *GRPCTestServer) ServerAddr() string {
	return s.Listener.Addr().String()
}

func (s *GRPCTestServer) TearDown() {
	if s.Listener != nil {
		s.Server.GracefulStop()
		s.Listener.Close()
	}
	if s.clientConn != nil {
		s.clientConn.Close()
	}
}
