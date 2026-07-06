type Server struct {
	proto.UnimplementedStreamServiceServer
}

func (s *Server) GetStream(req *proto.StreamReqData, stream grpc.ServerStreamingServer[proto.StreamResData]) error {
	fmt.Println("====== 服务端 流模式 ======")
	fmt.Println(req.GetData())
	for i := 0; i < 10; i++ {
		if err := stream.Send(&proto.StreamResData{
			Data: "Server: " + time.Now().Format("2006-01-02 15:04:05"),
		}); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (s *Server) PutStream(stream grpc.ClientStreamingServer[proto.StreamReqData, proto.StreamResData]) error {
	fmt.Println("====== 客户端 流模式 ======")
	count := 0
	for {
		if req, err := stream.Recv(); err == io.EOF {
			return stream.SendAndClose(&proto.StreamResData{
				Data: fmt.Sprintf("received %d messages", count),
			})
		} else if err != nil {
			return err
		} else {
			fmt.Println(req.GetData())
			count++
		}
	}
}

func (s *Server) AllStream(stream grpc.BidiStreamingServer[proto.StreamReqData, proto.StreamResData]) error {
	fmt.Println("======  双向  流模式 ======")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Println(req.GetData())
		if err := stream.Send(&proto.StreamResData{
			Data: "Server: " + time.Now().Format("2006-01-02 15:04:05"),
		}); err != nil {
			return err
		}
	}
}

func main() {
	server := grpc.NewServer()
	proto.RegisterStreamServiceServer(server, &Server{})
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = lis.Close(); err != nil {
			panic(err)
		}
	}()
	if err = server.Serve(lis); err != nil {
		panic(err)
	}
}
