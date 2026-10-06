func generateWithTool(ctx context.Context, base model.ToolCallingChatModel, noteTool tool.InvokableTool) (*schema.Message, error) {
	info, err := noteTool.Info(ctx)
	if err != nil {
		return nil, err
	}
	bound, err := base.WithTools([]*schema.ToolInfo{info})
	if err != nil {
		return nil, err
	}
	return bound.Generate(ctx, []*schema.Message{
		schema.UserMessage("请查找 P4 笔记的链接。"),
	})
}
