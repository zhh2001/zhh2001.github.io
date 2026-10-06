package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "example.com/grpcdemo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewStreamServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 服务端流：读到 EOF，才表示本次 RPC 成功结束。
	receiver, err := client.GetStream(ctx, &pb.StreamReqData{Data: "log"})
	if err != nil {
		return err
	}
	for {
		reply, err := receiver.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(reply.GetData())
	}

	// 客户端流：即使 Send 返回 EOF，也要获取服务端的最终状态。
	sender, err := client.PutStream(ctx)
	if err != nil {
		return err
	}
	for i := 1; i <= 3; i++ {
		err := sender.Send(&pb.StreamReqData{Data: fmt.Sprintf("item %d", i)})
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	reply, err := sender.CloseAndRecv()
	if err != nil {
		return err
	}
	fmt.Println(reply.GetData())
	return runBidi(ctx, client)
}

func runBidi(parent context.Context, client pb.StreamServiceClient) error {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	stream, err := client.AllStream(ctx)
	if err != nil {
		return err
	}
	sent := make(chan error, 1)
	go func() {
		for i := 1; i <= 3; i++ {
			err := stream.Send(&pb.StreamReqData{Data: fmt.Sprintf("chat %d", i)})
			if err == io.EOF {
				sent <- nil // 最终状态由接收方读取。
				return
			}
			if err != nil {
				sent <- err
				return
			}
		}
		sent <- stream.CloseSend()
	}()
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			return <-sent
		}
		if err != nil {
			cancel() // 让可能阻塞的发送方退出，再回收 goroutine。
			<-sent
			return err
		}
		fmt.Println(reply.GetData())
	}
}
