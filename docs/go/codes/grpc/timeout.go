ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
_, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Zhang"})
fmt.Println(status.Code(err)) // 调用上面的延迟服务端时，输出 DeadlineExceeded。
