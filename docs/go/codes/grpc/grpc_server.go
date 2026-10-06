package main

import (
	"context"
	"log"
	"net"

	pb "example.com/grpcdemo/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Msg: "Hello " + req.GetName()}, nil
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &Server{})
	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
