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

type State struct {
	History map[string]any
}

func getFunc(ctx context.Context) *State {
	return &State{
		History: make(map[string]any),
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	model, _ := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("MODEL"),
	})

	graph := compose.NewGraph[map[string]string, *schema.Message](compose.WithGenLocalState(getFunc))
	_ = graph.AddChatModelNode("model", model)

	lambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output map[string]string, err error) {

		// 内部使用 state
		_ = compose.ProcessState(ctx, func(ctx context.Context, s *State) error {
			s.History["tsundere_action"] = "我喜欢你"
			s.History["cute_action"] = "摸摸头"
			return nil
		})

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
		_ = compose.ProcessState(ctx, func(ctx context.Context, s *State) error {
			input["content"] = input["content"] + s.History["cute_action"].(string)
			return nil
		})

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

	tsundereStateToInput := func(ctx context.Context, in map[string]string, state *State) (output map[string]string, err error) {
		in["content"] = in["content"] + state.History["tsundere_action"].(string)
		return in, nil
	}

	_ = graph.AddLambdaNode("lambda", lambda)
	_ = graph.AddLambdaNode("tsundereLambda", tsundereLambda, compose.WithStatePreHandler(tsundereStateToInput))
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

	r, _ := graph.Compile(ctx)
	output, _ := r.Invoke(ctx, map[string]string{
		"role":    "tsundere",
		"content": "你好啊",
	})
	fmt.Println(output.Content)
}
