package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("MODEL"),
	})
	if err != nil {
		panic(err)
	}

	template := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个{role}"),
		&schema.Message{
			Role:    schema.User,
			Content: "请帮帮我，史瓦罗先生，帮我解决{task}",
		},
	)

	params := map[string]any{
		"role": "机器人史瓦罗先生",
		"task": "帮我刷星琼",
	}

	msg, err := template.Format(ctx, params)
	if err != nil {
		panic(err)
	}

	response, err := chatModel.Generate(ctx, msg)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Content)
}
