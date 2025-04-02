package signals

import (
	"os"
)

// ShutdownSignalHandler 表示终止信号处理器
type ShutdownSignalHandler interface {
	// SetupSignalHandler 设置并返回一个信号通道
	SetupSignalHandler() <-chan os.Signal

	// Cleanup 执行清理操作
	Cleanup() error
}
