func ProcessState[S any](ctx context.Context, handler func(context.Context, S) error) error {
	s, pMu, err := getState[S](ctx)
	if err != nil {
		return fmt.Errorf("get state from context fail: %w", err)
	}
	pMu.Lock()
	defer pMu.Unlock()
	return handler(ctx, s)
}