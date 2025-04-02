package graceful

import (
	"time"
)

// ShutdownManager 定义优雅终止管理器的接口
type ShutdownManager interface {
	// Start 启动监听终止信号
	Start() error

	// Stop 手动触发终止流程
	Stop() error

	// RegisterHook 注册在终止前需要执行的钩子函数
	RegisterHook(name string, hook ShutdownHook, timeout time.Duration) error

	// Wait 阻塞直到收到终止信号并完成所有钩子函数
	Wait() error

	// Done 返回一个channel，当收到终止信号时会关闭
	Done() <-chan struct{}
}

// NewManager 创建一个新的ShutdownManager
func NewManager(opts Options) ShutdownManager {
	// 实际创建会在平台特定文件中实现
	return newManager(opts)
}
