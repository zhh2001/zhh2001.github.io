ch := make(chan string)
timer := time.NewTimer(100 * time.Millisecond)
defer timer.Stop()

select {
case msg := <-ch:
	fmt.Println(msg)
case <-timer.C:
	fmt.Println("Timed out")
}
