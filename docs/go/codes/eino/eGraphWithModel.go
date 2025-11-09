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

	graph := compose.NewGraph[map[string]string, *schema.Message]()
	err = graph.AddChatModelNode("model", model)
	if err != nil {
		panic(err)
	}

	lambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output map[string]string, err error) {
		switch input["role"] {
		case "tsundere":
			return map[string]string{
				"role":    "傲娇",
				"content": input["content"],
			}, nil
		case "cute":
			return map[string]string{
				"role":    "可爱",
				"content": input["content"],
			}, nil
		default:
			return map[string]string{
				"role":    "user",
				"content": input["content"],
			}, nil
		}
	})

	tsundereLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个高冷傲娇的大小姐",
			},
			{
				Role:    schema.System,
				Content: input["content"],
			},
		}, nil
	})

	cuteLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个可爱的女高中生",
			},
			{
				Role:    schema.System,
				Content: input["content"],
			},
		}, nil
	})

	err = graph.AddLambdaNode("lambda", lambda)
	if err != nil {
		panic(err)
	}
	err = graph.AddLambdaNode("tsundereLambda", tsundereLambda)
	if err != nil {
		panic(err)
	}
	err = graph.AddLambdaNode("cuteLambda", cuteLambda)
	if err != nil {
		panic(err)
	}

	graphBranch := compose.NewGraphBranch(func(ctx context.Context, in map[string]string) (endNode string, err error) {
		switch in["role"] {
		case "傲娇":
			return "tsundereLambda", nil
		case "可爱":
			return "cuteLambda", nil
		default:
			return "tsundereLambda", nil
		}
	}, map[string]bool{
		"tsundereLambda": true,
		"cuteLambda":     true,
	})

	err = graph.AddBranch("lambda", graphBranch)
	if err != nil {
		panic(err)
	}
	err = graph.AddEdge(compose.START, "lambda")
	if err != nil {
		panic(err)
	}
	err = graph.AddEdge("tsundereLambda", "model")
	if err != nil {
		panic(err)
	}
	err = graph.AddEdge("cuteLambda", "model")
	if err != nil {
		panic(err)
	}
	err = graph.AddEdge("model", compose.END)
	if err != nil {
		panic(err)
	}

	r, err := graph.Compile(ctx)
	if err != nil {
		panic(err)
	}
	output, err := r.Invoke(ctx, map[string]string{
		"role":    "tsundere",
		"content": "你好啊",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(output.Content)
}
