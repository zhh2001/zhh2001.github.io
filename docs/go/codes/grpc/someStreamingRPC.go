func (s *Server) GetStream(req *pb.StreamReqData, stream grpc.ServerStreamingServer[pb.StreamResData]) error {
	md, _ := metadata.FromIncomingContext(stream.Context())
	fmt.Println("request-id:", md.Get("request-id"))
	if err := stream.SetHeader(metadata.Pairs("server", "grpcdemo")); err != nil {
		return err
	}
	stream.SetTrailer(metadata.Pairs("result", "done")) // 此方法没有返回值。
	return stream.Send(&pb.StreamResData{Data: req.GetData()})
}
