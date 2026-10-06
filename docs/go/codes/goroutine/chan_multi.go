package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(value int) {
			defer wg.Done()
			ch <- value
		}(i)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	for value := range ch {
		fmt.Println(value)
	}
}
