type Indexer interface {
	// Store stores the documents.
	Store(ctx context.Context, docs []*schema.Document, opts ...Option) (ids []string, err error) // invoke
}