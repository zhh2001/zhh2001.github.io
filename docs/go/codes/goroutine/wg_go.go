package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Go(func() {
			fmt.Println("Task:", i)
		})
	}
	wg.Wait()
	fmt.Println("All Done")
}
