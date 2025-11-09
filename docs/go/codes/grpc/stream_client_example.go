func main() {
	conn, err := grpc.NewClient(
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()

	client := proto.NewGreeterClient(conn)
	ctx := context.Background()

	fmt.Println("====== 服务端 流模式 ======")
	stream1, err := client.GetStream(ctx, &proto.StreamReqData{
		Data: "Client: " + time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		panic(err)
	}
	for {
		if res, err := stream1.Recv(); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			fmt.Println(res.GetData())
		}
	}

	fmt.Println("====== 客户端 流模式 ======")
	stream2, err := client.PutStream(ctx)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 10; i++ {
		if err = stream2.Send(&proto.StreamReqData{
			Data: "Client: " + time.Now().Format("2006-01-02 15:04:05"),
		}); err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}
	if err = stream2.CloseSend(); err != nil {
		panic(err)
	}

	fmt.Println("======  双向  流模式 ======")
	stream3, err := client.AllStream(ctx)
	wg := sync.WaitGroup{}

	// 发送协程
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if err = stream3.CloseSend(); err != nil {
				panic(err)
			}
		}()
		for i := 0; i < 10; i++ {
			if err = stream3.Send(&proto.StreamReqData{
				Data: "Client: " + time.Now().Format("2006-01-02 15:04:05"),
			}); err != nil {
				panic(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// 接收协程
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if res, err := stream3.Recv(); err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			} else {
				fmt.Println(res.GetData())
			}
		}
	}()

	wg.Wait()
}