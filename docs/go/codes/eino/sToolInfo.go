type ToolInfo struct {
	// 工具的唯一名称，用于清晰地表达其用途
	Name string
	// 用于告诉模型如何/何时/为什么使用这个工具
	// 可以在描述中包含少量示例
	Desc string
	// Extra is the extra information for the tool.
	Extra map[string]any

	// 工具接受的参数定义
	// 可以通过两种方式描述：
	//   1. 使用 ParameterInfo：schema.NewParamsOneOfByParams(params)
	//   2. 使用 OpenAPIV3：schema.NewParamsOneOfByOpenAPIV3(openAPIV3)
	// If is nil, signals that the tool does not need any input parameter
	*ParamsOneOf
}