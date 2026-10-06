type ToolInfo struct {
	Name  string
	Desc  string
	Extra map[string]any
	// 用 NewParamsOneOfByParams 或 NewParamsOneOfByJSONSchema 描述参数。
	*ParamsOneOf
}
