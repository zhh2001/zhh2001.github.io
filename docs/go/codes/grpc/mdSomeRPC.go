var header, trailer metadata.MD
r, err := client.SomeRPC(
    ctx,
    someRequest,
    grpc.Header(&header),
    grpc.Trailer(&trailer),
)