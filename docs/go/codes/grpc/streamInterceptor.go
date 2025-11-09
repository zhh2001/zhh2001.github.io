func interceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 拦截下来统计耗时
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)  // invoker 就是原先的调用逻辑
	fmt.Printf("耗时：%s", time.Since(start))
	return err
}

func main() {
	conn, err := grpc.NewClient(
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor),
	)
}