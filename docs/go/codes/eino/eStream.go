package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	apiKey, modelID := os.Getenv("ARK_API_KEY"), os.Getenv("MODEL")
	if apiKey == "" || modelID == "" {
		return fmt.Errorf("请设置 ARK_API_KEY 和 MODEL")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  apiKey,
		Model:   modelID,
		BaseURL: os.Getenv("ARK_BASE_URL"),
	})
	if err != nil {
		return err
	}
	input := []*schema.Message{
		schema.SystemMessage("你是一名技术助理，请用中文简要回答。"),
		schema.UserMessage("什么是软件定义网络？"),
	}
	reader, err := chatModel.Stream(ctx, input)
	if err != nil {
		return err
	}
	defer reader.Close()
	for {
		chunk, err := reader.Recv()
		if err == io.EOF {
			fmt.Println()
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Print(chunk.Content)
	}
}
