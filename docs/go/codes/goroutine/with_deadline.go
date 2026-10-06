deadline := time.Now().Add(10 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

actual, ok := ctx.Deadline()
fmt.Println(actual.Equal(deadline), ok) // true true
