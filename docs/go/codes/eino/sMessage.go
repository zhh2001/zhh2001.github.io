type Message struct {
	// Role 表示消息的角色（system/user/assistant/tool）
	Role RoleType `json:"role"`

	// // Content 是消息的文本内容
	Content string `json:"content"`

	// if MultiContent is not empty, use this instead of Content
	// if MultiContent is empty, use Content
	// MultiContent 是多模态内容，支持文本、图片、音频等
	// Deprecated: 已废弃，使用 UserInputMultiContent 替代
	MultiContent []ChatMessagePart `json:"multi_content,omitempty"`

	// UserInputMultiContent 用来存储用户输入的多模态数据，支持文本、图片、音频、视频、文件
	// 使用此字段时限制模型角色为 User
	UserInputMultiContent []MessageInputPart `json:"user_input_multi_content,omitempty"`

	// AAssistantGenMultiContent 用来承接模型输出的多模态数据，支持文本、图片、音频、视频
	// 使用此字段时限制模型角色为 Assistant
	AssistantGenMultiContent []MessageOutputPart `json:"assistant_output_multi_content,omitempty"`

	// Name 是消息的发送者名称
	Name string `json:"name,omitempty"`

	// ToolCalls 是 assistant 消息中的工具调用信息
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`

	// ToolCallID 是 tool 消息的工具调用 ID
	ToolCallID string `json:"tool_call_id,omitempty"`
	// only for ToolMessage
	ToolName string `json:"tool_name,omitempty"`

	// ResponseMeta 包含响应的元信息
	ResponseMeta *ResponseMeta `json:"response_meta,omitempty"`

	// ReasoningContent 是模型的推理过程，当模型返回推理内容时会包含该部分
	ReasoningContent string `json:"reasoning_content,omitempty"`

	// Extra 用于存储额外信息
	Extra map[string]any `json:"extra,omitempty"`
}