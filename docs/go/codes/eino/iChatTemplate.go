type ChatTemplate interface {
	Format(ctx context.Context, vs map[string]any, opts ...Option) ([]*schema.Message, error)
}