package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "example.com/grpcdemo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func logUnary(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("method=%s duration=%s code=%s", method, time.Since(start), status.Code(err))
	return err
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(logUnary))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	reply, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Zhang"})
	if err != nil {
		return err
	}
	fmt.Println(reply.GetMsg())
	return nil
}
