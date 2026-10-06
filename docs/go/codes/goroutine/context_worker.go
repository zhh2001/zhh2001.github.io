package main

import (
	"context"
	"fmt"
)

func worker(ctx context.Context, done chan<- struct{}) {
	defer close(done)
	<-ctx.Done()
	fmt.Println("Worker:", ctx.Err())
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{})
	go worker(ctx, done)

	cancel()
	<-done
	fmt.Println("All Done")
}
