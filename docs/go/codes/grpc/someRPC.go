func (s *server) SomeRPC(ctx context.Context, in *pb.SomeRequest) (*pb.SomeResponse, error) {
	header := metadata.Pairs("header-key", "val")
	if err := grpc.SetHeader(ctx, header); err != nil {
		return nil, err
	}
	trailer := metadata.Pairs("trailer-key", "val")
	if err := grpc.SetTrailer(ctx, trailer); err != nil {
		return nil, err
	}
	return &pb.SomeResponse{}, nil
}
