package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/compose"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	graph := compose.NewGraph[string, string]()
	classify := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		switch input {
		case "1":
			return "cat", nil
		case "2":
			return "dog", nil
		default:
			return "other", nil
		}
	})
	if err := graph.AddLambdaNode("classify", classify); err != nil {
		return err
	}
	for _, name := range []string{"cat", "dog", "other"} {
		lambda := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
			return map[string]string{"cat": "喵喵喵", "dog": "汪汪汪", "other": "你好"}[input], nil
		})
		if err := graph.AddLambdaNode(name, lambda); err != nil {
			return err
		}
		if err := graph.AddEdge(name, compose.END); err != nil {
			return err
		}
	}
	branch := compose.NewGraphBranch(func(ctx context.Context, input string) (string, error) {
		return input, nil
	}, map[string]bool{"cat": true, "dog": true, "other": true})
	if err := graph.AddBranch("classify", branch); err != nil {
		return err
	}
	if err := graph.AddEdge(compose.START, "classify"); err != nil {
		return err
	}
	runnable, err := graph.Compile(ctx)
	if err != nil {
		return err
	}
	for _, input := range []string{"1", "2", "3"} {
		output, err := runnable.Invoke(ctx, input)
		if err != nil {
			return err
		}
		fmt.Println(output)
	}
	return nil
}
