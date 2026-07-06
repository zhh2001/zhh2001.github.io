type Server struct {
	proto.UnimplementedGreeterServer // 为后续新增 RPC 保留兼容性
}

func (s *Server) SayHello(_ context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{Msg: "Hello " + req.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	proto.RegisterGreeterServer(server, &Server{})
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
