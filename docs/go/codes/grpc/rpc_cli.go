package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	client, err := rpc.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		return err
	}
	defer client.Close()
	var reply string
	if err := client.Call("HelloService.SayHello", "Zhang", &reply); err != nil {
		return err
	}
	fmt.Println(reply)
	return nil
}
