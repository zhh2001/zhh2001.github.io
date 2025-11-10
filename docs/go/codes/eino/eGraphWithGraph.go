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
	_ = godotenv.Load(".env")

	ctx := context.Background()
	model, _ := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("MODEL"),
	})

	graph := compose.NewGraph[map[string]string, *schema.Message]()
	_ = graph.AddChatModelNode("model", model)

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
			schema.SystemMessage("你是一个高冷傲娇的大小姐"),
			schema.UserMessage(input["content"]),
		}, nil
	})

	cuteLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		return []*schema.Message{
			schema.SystemMessage("你是一个可爱的女高中生"),
			schema.UserMessage(input["content"]),
		}, nil
	})

	_ = graph.AddLambdaNode("lambda", lambda)
	_ = graph.AddLambdaNode("tsundereLambda", tsundereLambda)
	_ = graph.AddLambdaNode("cuteLambda", cuteLambda)

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

	_ = graph.AddBranch("lambda", graphBranch)
	_ = graph.AddEdge(compose.START, "lambda")
	_ = graph.AddEdge("tsundereLambda", "model")
	_ = graph.AddEdge("cuteLambda", "model")
	_ = graph.AddEdge("model", compose.END)

	// 外层图
	outerGraph := compose.NewGraph[map[string]string, string]()
	outerLambda1 := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output map[string]string, err error) {
		return input, nil
	})
	outerLambda2 := compose.InvokableLambda(func(ctx context.Context, input *schema.Message) (output string, err error) {
		return input.Content + "\n+++ END +++", nil
	})

	_ = outerGraph.AddGraphNode("innerGraph", graph)
	_ = outerGraph.AddLambdaNode("outerLambda1", outerLambda1)
	_ = outerGraph.AddLambdaNode("outerLambda2", outerLambda2)
	_ = outerGraph.AddEdge(compose.START, "outerLambda1")
	_ = outerGraph.AddEdge("outerLambda1", "innerGraph")
	_ = outerGraph.AddEdge("innerGraph", "outerLambda2")
	_ = outerGraph.AddEdge("outerLambda2", compose.END)

	input := map[string]string{
		"role":    "cute",
		"content": "你好啊",
	}

	r, _ := outerGraph.Compile(ctx)
	output, _ := r.Invoke(ctx, input)
	fmt.Println(output)
}
