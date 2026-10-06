type Context interface {
	// 有截止时间时返回对应时间和 true
	Deadline() (deadline time.Time, ok bool)

	// 取消时关闭该通道，不可取消的上下文可以返回 nil
	Done() <-chan struct{}

	// 未取消时返回 nil，取消后返回 Canceled 或 DeadlineExceeded
	Err() error

	// 返回与 key 关联的值，如果没有则返回 nil
	Value(key any) any
}
