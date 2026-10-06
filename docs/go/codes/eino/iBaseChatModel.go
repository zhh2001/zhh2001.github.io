type messageType interface {
	*schema.Message | *schema.AgenticMessage
}

type BaseModel[M messageType] interface {
	Generate(ctx context.Context, input []M, opts ...Option) (M, error)
	Stream(ctx context.Context, input []M, opts ...Option) (*schema.StreamReader[M], error)
}

type BaseChatModel = BaseModel[*schema.Message]

type ToolCallingChatModel interface {
	BaseChatModel
	WithTools(tools []*schema.ToolInfo) (ToolCallingChatModel, error)
}
