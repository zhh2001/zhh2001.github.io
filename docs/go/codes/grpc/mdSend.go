// 创建新上下文，并添加一些元数据
ctx := metadata.AppendToOutgoingContext(ctx, "k1", "v1", "k1", "v2", "k2", "v3")

// 向上下文中添加更多元数据
ctx := metadata.AppendToOutgoingContext(ctx, "k3", "v4")

// 一元RPC
response, err := client.SomeRPC(ctx, someRequest)

// 流式RPC
stream, err := client.SomeStreamingRPC(ctx)