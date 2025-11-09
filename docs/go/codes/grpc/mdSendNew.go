// 创建新上下文，并添加一些元数据
md := metadata.Pairs("k1", "v1", "k1", "v2", "k2", "v3")
ctx := metadata.NewOutgoingContext(context.Background(), md)

// 向上下文中添加更多元数据
send, _ := metadata.FromOutgoingContext(ctx)
newMD := metadata.Pairs("k3", "v3")
ctx = metadata.NewOutgoingContext(ctx, metadata.Join(send, newMD))

// 一元RPC
response, err := client.SomeRPC(ctx, someRequest)

// 流式RPC
stream, err := client.SomeStreamingRPC(ctx)