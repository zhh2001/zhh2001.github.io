func (s *Server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return nil, status.Error(codes.NotFound, "user not found")
}
