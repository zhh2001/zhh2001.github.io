package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/callbacks"
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
	chain := compose.NewChain[string, string]()
	chain.AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		return strings.ToUpper(input), nil
	}), compose.WithNodeName("upper"))
	runnable, err := chain.Compile(ctx)
	if err != nil {
		return err
	}
	output, err := runnable.Invoke(ctx, "eino", compose.WithCallbacks(genCallback()))
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}

func genCallback() callbacks.Handler {
	return callbacks.NewHandlerBuilder().
		OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
			if info != nil {
				fmt.Printf("start: %s %s\n", info.Component, info.Name)
			}
			return ctx
		}).
		OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
			if info != nil {
				fmt.Printf("end: %s %s\n", info.Component, info.Name)
			}
			return ctx
		}).
		OnErrorFn(func(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
			fmt.Printf("error: %v\n", err)
			return ctx
		}).Build()
}
