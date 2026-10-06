func callWithMetadata(ctx context.Context, client pb.GreeterClient) error {
	var header, trailer metadata.MD
	reply, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Zhang"},
		grpc.Header(&header), grpc.Trailer(&trailer))
	fmt.Println("header:", header, "trailer:", trailer)
	if err != nil {
		return err
	}
	fmt.Println(reply.GetMsg())
	return nil
}
