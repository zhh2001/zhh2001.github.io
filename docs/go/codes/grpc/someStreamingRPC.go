func (s *server) SomeStreamingRPC(stream pb.Service_SomeStreamingRPCServer) error {
    header := metadata.Pairs("header-key", "val")
    if err := stream.SetHeader(header); err != nil {
        return err
    }
    trailer := metadata.Pairs("trailer-key", "val")
    stream.SetTrailer(trailer)
    return nil
}
