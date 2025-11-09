ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
_, err := client.SayHello(ctx, &proto.HelloRequest{})
if err != nil {
	if st, ok := status.FromError(err); ok {
		fmt.Println(st.Code())     // DeadlineExceeded
		fmt.Println(st.Message())  // context deadline exceeded
	}
}