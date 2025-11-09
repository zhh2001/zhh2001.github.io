func sayHello() {
    fmt.Println("Hello from goroutine!")
}

func main() {
    go sayHello()  // 启动一个新的协程
    fmt.Println("Hello from main!")
    time.Sleep(1 * time.Second)  // 等待协程完成
}