func (s *server) SomeStreamingRPC(stream pb.Service_SomeStreamingRPCServer) error {
    header  := metadata.Pairs("header-key", "val")
    stream.SetHeader(header)
    trailer := metadata.Pairs("trailer-key", "val")
    stream.SetTrailer(trailer)
}