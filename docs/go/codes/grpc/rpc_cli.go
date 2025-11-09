// 连接到 RPC 服务
client, _ := rpc.Dial("tcp", "127.0.0.1:1234")

// 调用远程方法
var reply string
_ = client.Call("HelloService.SayHello", "zhang", &reply)
fmt.Println(reply)  // Hello zhang