ch := make(chan int, 2)

ch <- 1
ch <- 2
// ch <- 3  // 如果取消注释，这里会阻塞，因为缓冲区已满

fmt.Println(<-ch)  // 输出 1
fmt.Println(<-ch)  // 输出 2