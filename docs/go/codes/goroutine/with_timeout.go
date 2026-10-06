ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

_, ok := ctx.Deadline()
fmt.Println(ok) // true
