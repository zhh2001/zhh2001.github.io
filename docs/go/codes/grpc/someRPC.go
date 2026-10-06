func (s *Server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Println("request-id:", md.Get("request-id"))
	if err := grpc.SetHeader(ctx, metadata.Pairs("server", "grpcdemo")); err != nil {
		return nil, err
	}
	if err := grpc.SetTrailer(ctx, metadata.Pairs("result", "done")); err != nil {
		return nil, err
	}
	return &pb.HelloResponse{Msg: "Hello " + req.GetName()}, nil
}
