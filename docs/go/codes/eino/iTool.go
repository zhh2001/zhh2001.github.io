// BaseTool 基础工具接口，提供工具信息
type BaseTool interface {
	Info(ctx context.Context) (*schema.ToolInfo, error)
}

// InvokableTool 支持同步调用的工具接口
type InvokableTool interface {
	BaseTool

	// InvokableRun call function with arguments in JSON format
	InvokableRun(ctx context.Context, argumentsInJSON string, opts ...Option) (string, error)
}

// StreamableTool 支持流式输出的工具接口
type StreamableTool interface {
	BaseTool

	StreamableRun(ctx context.Context, argumentsInJSON string, opts ...Option) (*schema.StreamReader[string], error)
}