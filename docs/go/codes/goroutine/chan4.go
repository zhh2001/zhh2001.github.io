ch := make(chan int, 2)
ch <- 1
ch <- 2
close(ch)

for v := range ch {
	fmt.Println(v)
}

v, ok := <-ch
fmt.Println(v, ok) // 0 false
