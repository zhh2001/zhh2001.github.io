// 生产者函数，只发送数据
func producer(sendCh chan<- int) {
	for i := 0; i < 5; i++ {
		sendCh <- i
		fmt.Println("Produced:", i)
	}
	close(sendCh)
}

// 消费者函数，只接收数据
func consumer(receiveCh <-chan int) {
	for num := range receiveCh {
		fmt.Println("Consumed:", num)
	}
}

func main() {
	ch := make(chan int)

	go producer(ch)
	consumer(ch)
}