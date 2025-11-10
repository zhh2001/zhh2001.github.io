import "github.com/cloudwego/eino/callbacks"

func genCallback() callbacks.Handler {
	handlerBuilder := callbacks.NewHandlerBuilder()
	handlerBuilder.OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
		fmt.Printf("节点：%s，输出：%#v\n", info.Component, input)
		return ctx
	})
	handlerBuilder.OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
		fmt.Printf("节点：%s，输出：%+v\n", info.Component, output)
		return ctx
	})
	return handlerBuilder.Build()
}

func main() {
	// ...
	r, _ := graph.Compile(ctx)
	output, _ := r.Invoke(ctx, input, compose.WithCallbacks(genCallback()))
	fmt.Println(output.Content)
}
