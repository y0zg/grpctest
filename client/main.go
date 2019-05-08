package main

import (
	"crypto/tls"
	"context"
	"fmt"
	"log"

	pb "github.com/zenoss/grpctest/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// var addr = "35.244.172.248:443"
var addr = "jpl.zenoss.io:443"
func main() {
	// opt := grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))

	// conn, err := grpc.Dial(addr, opt)
	// dialOpts := []grpc.DialOption{
	// 	grpc.WithBlock(),
	// 	grpc.WithInsecure(),
	// }

	tlsCreds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	opt := grpc.WithTransportCredentials(tlsCreds)

	conn, err := grpc.Dial(addr, opt)
	if err != nil {
		log.Fatalf("You do not deserve a connection: %v", err)
	}
	defer conn.Close()
	client := pb.NewMathServiceClient(conn)

	resp, err := client.Square(context.Background(), &pb.Request{Value: 20})
	if err != nil {
		log.Fatalf("A bad request: %v", err)
	}
	fmt.Println("Received value:", resp.Value)

	resp, err = client.Random(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("A bad request: %v", err)
	}
	fmt.Println("Received Random value:", resp.Value)

}
