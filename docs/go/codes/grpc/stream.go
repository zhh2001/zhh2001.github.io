type ChatServer struct {
	pb.UnimplementedChatServer
}

func (s *ChatServer) Chat(stream grpc.BidiStreamingServer[pb.ChatMessage, pb.ChatMessage]) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&pb.ChatMessage{Text: "echo: " + msg.GetText()}); err != nil {
			return err
		}
	}
}
