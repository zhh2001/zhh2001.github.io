type UploadServer struct {
	pb.UnimplementedUploaderServer
}

func (s *UploadServer) UploadFile(stream grpc.ClientStreamingServer[pb.FileChunk, pb.FileSummary]) error {
	var size uint64
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.FileSummary{Size: size})
		}
		if err != nil {
			return err
		}
		size += uint64(len(chunk.GetData()))
	}
}
