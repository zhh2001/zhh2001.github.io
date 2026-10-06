package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/compose"
	"log"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	inner := compose.NewGraph[string, string]()
	upper := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		return strings.ToUpper(strings.TrimSpace(input)), nil
	})
	if err := inner.AddLambdaNode("upper", upper); err != nil {
		return err
	}
	for _, edge := range [][2]string{{compose.START, "upper"}, {"upper", compose.END}} {
		if err := inner.AddEdge(edge[0], edge[1]); err != nil {
			return err
		}
	}
	outer := compose.NewGraph[string, string]()
	if err := outer.AddGraphNode("inner", inner); err != nil {
		return err
	}
	decorate := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		return "result: " + input, nil
	})
	if err := outer.AddLambdaNode("decorate", decorate); err != nil {
		return err
	}
	for _, edge := range [][2]string{{compose.START, "inner"}, {"inner", "decorate"}, {"decorate", compose.END}} {
		if err := outer.AddEdge(edge[0], edge[1]); err != nil {
			return err
		}
	}
	runnable, err := outer.Compile(ctx)
	if err != nil {
		return err
	}
	output, err := runnable.Invoke(ctx, " eino ")
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}
