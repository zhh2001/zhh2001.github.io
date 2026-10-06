// 公开 API 签名，省略框架内部的状态查找和加锁实现。
func ProcessState[S any](ctx context.Context, handler func(context.Context, S) error) error
