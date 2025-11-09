// 服务端
func (s *Server) Chat(stream pb.Service_ChatServer) error {
    for {
        msg, err := stream.Recv()
        if err != nil {
            return err
        }
        // 处理消息并回复
        stream.Send(&pb.ChatResponse{...})
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