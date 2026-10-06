package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "example.com/grpcdemo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Msg: "Hello " + req.GetName()}, nil
}

func logUnary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	reply, err := handler(ctx, req)
	log.Printf("method=%s duration=%s code=%s", info.FullMethod, time.Since(start), status.Code(err))
	return reply, err
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(logUnary))
	pb.RegisterGreeterServer(server, &Server{})
	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
