type HelloService struct{}

func (s *HelloService) SayHello(request string, reply *string) error {
	*reply = "Hello " + request
	return nil
}

func main() {
	if err := rpc.RegisterName("HelloService", new(HelloService)); err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
