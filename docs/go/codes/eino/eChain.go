package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
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
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一名{role}，请用中文回答。"),
		schema.MessagesPlaceholder("history", true),
		schema.UserMessage("请解释{topic}。"),
	)
	chain := compose.NewChain[map[string]any, *schema.Message]()
	chain.AppendChatTemplate(template).AppendChatModel(chatModel)
	runnable, err := chain.Compile(ctx)
	if err != nil {
		return err
	}
	output, err := runnable.Invoke(ctx, map[string]any{
		"role": "网络技术助理", "topic": "软件定义网络",
	})
	if err != nil {
		return err
	}
	fmt.Println(output.Content)
	return nil
}
