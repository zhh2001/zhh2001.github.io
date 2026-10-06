type Options struct {
	Index    *string
	SubIndex *string
	TopK     *int
	// 阈值的含义和支持情况由检索器实现决定。
	ScoreThreshold *float64
	Embedding      embedding.Embedder
	// 后端特有的过滤或查询参数，不限定为某一种检索器。
	DSLInfo map[string]any
}
