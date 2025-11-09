func main() {
	conn, _ := grpc.NewClient(
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	defer conn.Close()
	client := proto.NewGreeterClient(conn)

	// 调用服务
	response, _ := client.SayHello(
		context.Background(),
		&proto.HelloRequest{Name: "Zhang"},
	)
	fmt.Println(response.GetMsg())
}