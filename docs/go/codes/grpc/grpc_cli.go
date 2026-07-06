func main() {
	conn, err := grpc.NewClient(
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := proto.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	response, err := client.SayHello(
		ctx,
		&proto.HelloRequest{Name: "Zhang"},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.GetMsg())
}
