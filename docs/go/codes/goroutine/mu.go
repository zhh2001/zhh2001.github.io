package main

import (
	"fmt"
	"sync"
)

var counter int
var mu sync.Mutex
var wg sync.WaitGroup

func count() {
	defer wg.Done()
	mu.Lock()
	defer mu.Unlock()
	counter++
}

func main() {
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go count()
	}
	wg.Wait()
	fmt.Println(counter) // 1000
}
