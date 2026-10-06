type LogServer struct {
	pb.UnimplementedLoggerServer
}

func (s *LogServer) StreamLogs(req *pb.LogRequest, stream grpc.ServerStreamingServer[pb.LogResponse]) error {
	for i := 1; i <= 3; i++ {
		reply := &pb.LogResponse{Line: fmt.Sprintf("%s: log %d", req.GetQuery(), i)}
		if err := stream.Send(reply); err != nil {
			return err
		}
	}
	return nil
}
