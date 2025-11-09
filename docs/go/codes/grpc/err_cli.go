if _, err := client.SayHello(context.Background(), nil); err != nil {
	if st, ok := status.FromError(err); ok {
		fmt.Println(st.Code(), st.Message())
	}
}