// 创建一个传递 int 类型的无缓冲 Channel
ch := make(chan int)

// 创建一个传递 int 类型的有缓冲 Channel，缓冲区大小为 10
chBuffered := make(chan int, 10)