ctx, cancel := context.WithCancel(context.Background())
defer cancel()  // 确保在不需要时取消上下文，释放资源