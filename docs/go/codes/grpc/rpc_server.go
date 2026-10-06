package main

import (
	"log"
	"net"
	"net/rpc"
)

type HelloService struct{}

func (s *HelloService) SayHello(name string, reply *string) error {
	*reply = "Hello " + name
	return nil
}

func main() {
	server := rpc.NewServer()
	if err := server.RegisterName("HelloService", new(HelloService)); err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			return
		}
		go server.ServeConn(conn)
	}
}
