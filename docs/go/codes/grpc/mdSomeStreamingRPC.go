func readWithMetadata(parent context.Context, client pb.StreamServiceClient) error {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	stream, err := client.GetStream(ctx, &pb.StreamReqData{Data: "log"})
	if err != nil {
		return err
	}
	header, err := stream.Header()
	if err != nil {
		return err
	}
	fmt.Println("header:", header)
	for {
		reply, err := stream.Recv()
		if err != nil {
			fmt.Println("trailer:", stream.Trailer())
			if err == io.EOF {
				return nil
			}
			return err
		}
		fmt.Println(reply.GetData())
	}
}
