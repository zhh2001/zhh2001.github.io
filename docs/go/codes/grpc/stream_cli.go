// 服务端
func (s *Server) UploadFile(stream grpc.ClientStreamingServer[pb.FileChunk, pb.FileSummary]) error {
    for {
        chunk, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&pb.FileSummary{...})
        }
        if err != nil {
            return err
        }
        // 处理 chunk
    }
}

// 客户端
stream, err := client.UploadFile(ctx)
for _, chunk := range chunks {
    if err := stream.Send(chunk); err != nil {
        break
    }
}
summary, err := stream.CloseAndRecv()
