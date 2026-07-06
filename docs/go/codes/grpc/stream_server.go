// 服务端
func (s *Server) StreamLogs(req *pb.LogRequest, stream grpc.ServerStreamingServer[pb.LogResponse]) error {
    for {
        if err := stream.Send(&pb.LogResponse{...}); err != nil {
            return err
        }
        time.Sleep(1 * time.Second)
    }
}

// 客户端
stream, err := client.StreamLogs(ctx, &pb.LogRequest{...})
for {
    res, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        return err
    }
    log.Println(res)
}
