func (s *Server) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	err := status.New(codes.NotFound, "Not Found").Err()
	return nil, err
}