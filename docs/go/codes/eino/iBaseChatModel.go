type BaseChatModel interface {
	Generate(ctx context.Context, input []*schema.Message, opts ...Option) (*schema.Message, error)
	Stream(ctx context.Context, input []*schema.Message, opts ...Option) (
		*schema.StreamReader[*schema.Message], error)
}

type ToolCallingChatModel interface {
	BaseChatModel

	// WithTools returns a new ToolCallingChatModel instance with the specified tools bound.
	// This method does not modify the current instance, making it safer for concurrent use.
	WithTools(tools []*schema.ToolInfo) (ToolCallingChatModel, error)
}