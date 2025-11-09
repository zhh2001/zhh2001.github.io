// Retriever is the interface for retriever.
// It is used to retrieve documents from a source.
//
// e.g.
//
//		retriever, err := redis.NewRetriever(ctx, &redis.RetrieverConfig{})
//		if err != nil {...}
//		docs, err := retriever.Retrieve(ctx, "query") // <= using directly
//		docs, err := retriever.Retrieve(ctx, "query", retriever.WithTopK(3)) // <= using options
//
//	 	graph := compose.NewGraph[inputType, outputType](compose.RunTypeDAG)
//		graph.AddRetrieverNode("retriever_node_key", retriever) // <= using in graph
type Retriever interface {
	Retrieve(ctx context.Context, query string, opts ...Option) ([]*schema.Document, error)
}