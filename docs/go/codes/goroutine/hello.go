package main

import "fmt"

func sayHello() {
	fmt.Println("Hello from goroutine!")
}

func main() {
	done := make(chan struct{})
	go func() {
		defer close(done)
		sayHello()
	}()
	fmt.Println("Hello from main!")
	<-done
}
