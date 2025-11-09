func interceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	// 拦截下来统计耗时
	start := time.Now()
	resp, err = handler(ctx, req)  // handler 就是原来正常的处理逻辑
	fmt.Printf("耗时：%s", time.Since(start))
	return resp, err
}

func main() {
	opt := grpc.UnaryInterceptor(interceptor)
	server := grpc.NewServer(opt)
}