type Context interface {
	// 返回上下文的截止时间
	Deadline() (deadline time.Time, ok bool)

	// 返回一个通道，当上下文被取消或到达截止时间时关闭
	Done() <-chan struct{}

	// 返回上下文被取消的原因
	Err() error

	// 返回与 key 关联的值，如果没有则返回 nil
	Value(key interface{}) interface{}
}