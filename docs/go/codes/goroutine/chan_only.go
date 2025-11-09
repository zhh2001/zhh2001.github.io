// 创建一个双向的整数通道
ch := make(chan int)

// sendOnly 只能用于发送数据
var sendOnly chan<- int = ch

// receiveOnly 只能用于接收数据
var receiveOnly <-chan int = ch