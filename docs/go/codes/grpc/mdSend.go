ctx := metadata.AppendToOutgoingContext(context.Background(), "request-id", "req-001")
ctx = metadata.AppendToOutgoingContext(ctx, "tag", "a", "tag", "b")
md, _ := metadata.FromOutgoingContext(ctx)
fmt.Println(md.Get("request-id"), md.Get("tag")) // [req-001] [a b]
// 将 ctx 传给 client.SayHello 或流式方法，元数据才会随 RPC 发出。
