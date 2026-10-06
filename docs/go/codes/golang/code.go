//go:generate stringer -type Code -linecomment

package code

type Code int64

const (
	OK            Code = 0 // OK
	InvalidParams Code = 1 // 参数错误
	Timeout       Code = 2 // 超时
)
