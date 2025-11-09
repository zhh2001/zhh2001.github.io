type Server struct {
	proto.UnimplementedGreeterServer // 必须嵌入
}

func (s Server) GetStream(req *proto.StreamReqData, stream grpc.ServerStreamingServer[proto.StreamResData]) error {
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

func (s Server) PutStream(stream grpc.ClientStreamingServer[proto.StreamReqData, proto.StreamResData]) error {
	fmt.Println("====== 客户端 流模式 ======")
	for {
		if req, err := stream.Recv(); err == io.EOF {
			break
		} else if err != nil {
			return err
		} else {
			fmt.Println(req.GetData())
		}
	}
	return nil
}

func (s Server) AllStream(stream grpc.BidiStreamingServer[proto.StreamReqData, proto.StreamResData]) error {
	fmt.Println("======  双向  流模式 ======")
	wg := sync.WaitGroup{}

	// 接收协程
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if req, err := stream.Recv(); err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			} else {
				fmt.Println(req.GetData())
			}
		}
	}()

	// 发送协程
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			if err := stream.Send(&proto.StreamResData{
				Data: "Server: " + time.Now().Format("2006-01-02 15:04:05"),
			}); err != nil {
				panic(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	wg.Wait()
	return nil
}

func main() {
	server := grpc.NewServer()
	s := &Server{}
	proto.RegisterGreeterServer(server, s)
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