ctx := context.Background()
fmt.Println(ctx.Err(), ctx.Done() == nil) // <nil> true
