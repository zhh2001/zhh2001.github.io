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
	graph := compose.NewGraph[string, string](compose.WithGenLocalState(func(ctx context.Context) *State {
		return &State{}
	}))
	count := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		err := compose.ProcessState(ctx, func(ctx context.Context, state *State) error {
			state.Count++
			return nil
		})
		return input, err
	})
	format := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		return input, nil
	})
	pre := func(ctx context.Context, input string, state *State) (string, error) {
		return fmt.Sprintf("%s: count=%d", input, state.Count), nil
	}
	if err := graph.AddLambdaNode("count", count); err != nil {
		return err
	}
	if err := graph.AddLambdaNode("format", format, compose.WithStatePreHandler(pre)); err != nil {
		return err
	}
	for _, edge := range [][2]string{{compose.START, "count"}, {"count", "format"}, {"format", compose.END}} {
		if err := graph.AddEdge(edge[0], edge[1]); err != nil {
			return err
		}
	}
	runnable, err := graph.Compile(ctx)
	if err != nil {
		return err
	}
	for i := 0; i < 2; i++ {
		output, err := runnable.Invoke(ctx, "request")
		if err != nil {
			return err
		}
		fmt.Println(output)
	}
	return nil
}

type State struct{ Count int }
