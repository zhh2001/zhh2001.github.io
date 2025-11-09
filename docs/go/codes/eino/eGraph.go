package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
)

func main() {
	ctx := context.Background()

	lambda0 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		switch input {
		case "1":
			return "猫", nil
		case "2":
			return "狗", nil
		default:
			return "人", nil
		}
	})

	lambda1 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		return "喵喵喵", nil
	})

	lambda2 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		return "汪汪汪", nil
	})

	lambda3 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		return "嘤嘤嘤", nil
	})

	graph := compose.NewGraph[string, string]()
	err := graph.AddLambdaNode("lambda0", lambda0)
	if err != nil {
		panic(err)
	}

	err = graph.AddLambdaNode("lambda1", lambda1)
	if err != nil {
		panic(err)
	}

	err = graph.AddLambdaNode("lambda2", lambda2)
	if err != nil {
		panic(err)
	}

	err = graph.AddLambdaNode("lambda3", lambda3)
	if err != nil {
		panic(err)
	}

	err = graph.AddBranch("lambda0", compose.NewGraphBranch(func(ctx context.Context, in string) (endNode string, err error) {
		switch in {
		case "猫":
			return "lambda1", nil
		case "狗":
			return "lambda2", nil
		case "人":
			return "lambda3", nil
		default:
			return compose.END, nil
		}
	}, map[string]bool{
		"lambda1": true,
		"lambda2": true,
		"lambda3": true,
	}))
	if err != nil {
		panic(err)
	}

	err = graph.AddEdge(compose.START, "lambda0")
	if err != nil {
		panic(err)
	}

	err = graph.AddEdge("lambda1", compose.END)
	if err != nil {
		panic(err)
	}

	err = graph.AddEdge("lambda2", compose.END)
	if err != nil {
		panic(err)
	}

	err = graph.AddEdge("lambda3", compose.END)
	if err != nil {
		panic(err)
	}

	r, err := graph.Compile(ctx)
	if err != nil {
		panic(err)
	}

	output, err := r.Invoke(ctx, "3")
	if err != nil {
		panic(err)
	}
	fmt.Println(output)
}
