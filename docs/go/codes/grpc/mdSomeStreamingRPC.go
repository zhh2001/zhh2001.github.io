stream, err := client.SomeStreamingRPC(ctx)
if err != nil {
    return err
}
header, err := stream.Header()
if err != nil {
    return err
}
for {
    _, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        return err
    }
}
trailer := stream.Trailer()
