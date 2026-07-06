client, err := rpc.Dial("tcp", "127.0.0.1:1234")
if err != nil {
	log.Fatal(err)
}
defer client.Close()

var reply string
if err := client.Call("HelloService.SayHello", "zhang", &reply); err != nil {
	log.Fatal(err)
}
fmt.Println(reply) // Hello zhang
