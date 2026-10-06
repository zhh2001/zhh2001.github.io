ctx := metadata.AppendToOutgoingContext(context.Background(), "tag", "old")
ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("tag", "new"))
md, _ := metadata.FromOutgoingContext(ctx)
fmt.Println(md.Get("tag")) // [new]

// 需要保留已有值时，先合并，再创建新的 outgoing context。
md = metadata.Join(md, metadata.Pairs("tag", "extra"))
ctx = metadata.NewOutgoingContext(ctx, md)
