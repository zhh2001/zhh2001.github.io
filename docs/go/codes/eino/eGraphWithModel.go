package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/ark"
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
	graph := compose.NewGraph[map[string]string, *schema.Message]()
	choose := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (map[string]string, error) {
		if input["style"] != "brief" && input["style"] != "detail" {
			return nil, fmt.Errorf("style 应为 brief 或 detail")
		}
		return input, nil
	})
	if err := graph.AddLambdaNode("choose", choose); err != nil {
		return err
	}
	for _, style := range []string{"brief", "detail"} {
		prepare := compose.InvokableLambda(func(ctx context.Context, input map[string]string) ([]*schema.Message, error) {
			instruction := "请简要解释技术概念。"
			if input["style"] == "detail" {
				instruction = "请分步骤解释技术概念，并给出一个例子。"
			}
			return []*schema.Message{schema.SystemMessage(instruction), schema.UserMessage(input["content"])}, nil
		})
		if err := graph.AddLambdaNode(style, prepare); err != nil {
			return err
		}
	}
	if err := graph.AddChatModelNode("model", chatModel); err != nil {
		return err
	}
	branch := compose.NewGraphBranch(func(ctx context.Context, input map[string]string) (string, error) {
		return input["style"], nil
	}, map[string]bool{"brief": true, "detail": true})
	if err := graph.AddBranch("choose", branch); err != nil {
		return err
	}
	for _, edge := range [][2]string{{compose.START, "choose"}, {"brief", "model"}, {"detail", "model"}, {"model", compose.END}} {
		if err := graph.AddEdge(edge[0], edge[1]); err != nil {
			return err
		}
	}
	runnable, err := graph.Compile(ctx)
	if err != nil {
		return err
	}
	output, err := runnable.Invoke(ctx, map[string]string{"style": "brief", "content": "什么是 P4？"})
	if err != nil {
		return err
	}
	fmt.Println(output.Content)
	return nil
}
