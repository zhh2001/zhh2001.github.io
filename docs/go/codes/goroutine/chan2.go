ch := make(chan string)

go func() {
    time.Sleep(2 * time.Second)
    ch <- "Zhang"
}()

msg := <-ch
fmt.Println(msg)