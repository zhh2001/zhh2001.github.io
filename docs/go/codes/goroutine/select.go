ch1 := make(chan string, 1)
ch2 := make(chan string, 1)
ch1 <- "ch1"
ch2 <- "ch2"

select {
case msg1 := <-ch1:
	fmt.Println(msg1)
case msg2 := <-ch2:
	fmt.Println(msg2)
}
