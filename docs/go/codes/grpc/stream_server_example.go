package main

import (
	"fmt"
	"io"
	"log"
	"net"

	pb "example.com/grpcdemo/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedStreamServiceServer
}

func (s *Server) GetStream(req *pb.StreamReqData, stream grpc.ServerStreamingServer[pb.StreamResData]) error {
	for i := 1; i <= 3; i++ {
		reply := &pb.StreamResData{Data: fmt.Sprintf("%s #%d", req.GetData(), i)}
		if err := stream.Send(reply); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) PutStream(stream grpc.ClientStreamingServer[pb.StreamReqData, pb.StreamResData]) error {
	count := 0
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamResData{Data: fmt.Sprintf("received %d messages", count)})
		}
		if err != nil {
			return err
		}
		count++
	}
}

func (s *Server) AllStream(stream grpc.BidiStreamingServer[pb.StreamReqData, pb.StreamResData]) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&pb.StreamResData{Data: "echo: " + req.GetData()}); err != nil {
			return err
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	pb.RegisterStreamServiceServer(server, &Server{})
	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
