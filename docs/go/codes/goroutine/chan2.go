ch := make(chan string)

go func() {
	ch <- "Zhang"
}()

msg := <-ch
fmt.Println(msg) // Zhang
