func executeToolCall(ctx context.Context) ([]*schema.Message, error) {
	node, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{CreateTool()},
	})
	if err != nil {
		return nil, err
	}
	message := schema.AssistantMessage("", []schema.ToolCall{{
		ID:       "call_1",
		Type:     "function",
		Function: schema.FunctionCall{Name: "get_note", Arguments: `{"name":"P4"}`},
	}})
	return node.Invoke(ctx, message)
}
