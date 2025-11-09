func (s *server) SomeRPC(ctx context.Context, in *pb.someRequest) (*pb.someResponse, error) {
	header := metadata.Pairs("header-key", "val")
	grpc.SetHeader(ctx, header)
	trailer := metadata.Pairs("trailer-key", "val")
	grpc.SetTrailer(ctx, trailer)
}