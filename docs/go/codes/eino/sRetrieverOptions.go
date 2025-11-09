type Options struct {
	// Index 是检索器使用的索引，不同检索器中的索引可能有不同含义
	Index *string
	// SubIndex 是检索器使用的子索引，不同检索器中的子索引可能有不同含义
	SubIndex *string
	// TopK 是检索的文档数量上限
	TopK *int
	// ScoreThreshold 是文档相似度的阈值，例如 0.5 表示文档的相似度分数必须大于 0.5
	ScoreThreshold *float64
	// Embedding 是用于生成查询向量的组件
	Embedding embedding.Embedder

	// DSLInfo 是用于检索的 DSL 信息，仅在 viking 类型的检索器中使用
	DSLInfo map[string]interface{}
}