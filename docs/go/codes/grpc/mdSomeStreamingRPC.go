stream, err := client.SomeStreamingRPC(ctx)
header, err := stream.Header()
trailer     := stream.Trailer()