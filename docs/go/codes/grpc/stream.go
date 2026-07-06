// 服务端
func (s *Server) Chat(stream grpc.BidiStreamingServer[pb.ChatMessage, pb.ChatMessage]) error {
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
        // 处理消息并回复
        if err := stream.Send(&pb.ChatMessage{...}); err != nil {
            return err
        }
    }
}

// 客户端
stream, err := client.Chat(ctx)
// 发送协程
go func() {
    for {
        stream.Send(&pb.ChatMessage{...})
    }
}()
// 接收协程
go func() {
    for {
        res, err := stream.Recv()
        // 处理响应
    }
}()
