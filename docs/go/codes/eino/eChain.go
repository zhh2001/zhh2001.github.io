package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("MODEL"),
	})
	if err != nil {
		panic(err)
	}

	lambda := compose.InvokableLambda(func(ctx context.Context, input string) (output []*schema.Message, err error) {
		meow := input + "回答结尾加上‘喵’"
		output = []*schema.Message{
			{
				Role:    schema.User,
				Content: meow,
			},
		}
		return output, nil
	})

	chain := compose.NewChain[string, *schema.Message]()
	chain.AppendLambda(lambda).AppendChatModel(model)
	r, err := chain.Compile(ctx)
	if err != nil {
		panic(err)
	}

	answer, err := r.Invoke(ctx, "你好，可以告诉我你的名字吗")
	if err != nil {
		panic(err)
	}
	fmt.Println(answer.Content)
}
