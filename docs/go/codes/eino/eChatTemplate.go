package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一名{role}，请用中文回答。"),
		schema.MessagesPlaceholder("history", true),
		schema.UserMessage("请解释{topic}。"),
	)
	params := map[string]any{
		"role":  "网络技术助理",
		"topic": "软件定义网络",
		"history": []*schema.Message{
			schema.UserMessage("我了解基本的 TCP/IP 协议。"),
			schema.AssistantMessage("接下来可以学习控制平面和数据平面的分工。", nil),
		},
	}
	messages, err := template.Format(ctx, params)
	if err != nil {
		return err
	}
	for _, message := range messages {
		fmt.Printf("%s: %s\n", message.Role, message.Content)
	}
	return nil
}
