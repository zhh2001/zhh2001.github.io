ctx, cancel := context.WithCancel(context.Background())
defer cancel()

fmt.Println(ctx.Err()) // <nil>
cancel()
<-ctx.Done()
fmt.Println(ctx.Err()) // context canceled
