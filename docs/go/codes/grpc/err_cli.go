func callAndCheck(ctx context.Context, client pb.GreeterClient) error {
	_, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Zhang"})
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	fmt.Println(st.Code(), st.Message())
	if st.Code() == codes.NotFound {
		fmt.Println("用户不存在")
	}
	return err
}
