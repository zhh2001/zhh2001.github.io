type Server struct {
	proto.UnimplementedGreeterServer // 必须嵌入
}

func (s *Server) SayHello(_ context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{Msg: "Hello " + req.Name}, nil
}

func main() {
	server := grpc.NewServer()
	proto.RegisterGreeterServer(server, &Server{})
	lis, _ := net.Listen("tcp", ":8080")
	server.Serve(lis)
}