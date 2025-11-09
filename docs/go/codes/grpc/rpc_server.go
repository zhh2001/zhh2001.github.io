type HelloService struct{}

func (s *HelloService) SayHello(request string, reply *string) error {
	*reply = "Hello " + request
	return nil
}

func main() {
	_ = rpc.RegisterName("HelloService", new(HelloService))

	// 监听 TCP 端口
	listener, _ := net.Listen("tcp", ":1234")

	// 接受连接并处理请求
	for {
		conn, _ := listener.Accept()
		go rpc.ServeConn(conn)
	}
}