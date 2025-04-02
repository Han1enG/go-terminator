package graceful

import (
	"context"
	"time"
)

// ShutdownHook 是终止钩子函数的接口
type ShutdownHook interface {
	// Execute 在终止过程中执行的操作
	Execute(ctx context.Context) error
}

// HookFunc 函数类型实现ShutdownHook接口
type HookFunc func(ctx context.Context) error

// Execute 实现ShutdownHook接口
func (f HookFunc) Execute(ctx context.Context) error {
	return f(ctx)
}

// hookEntry 内部使用的钩子条目
type hookEntry struct {
	hook    ShutdownHook
	timeout time.Duration
}
