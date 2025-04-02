//go:build linux
// +build linux

package signals

import (
	"os"
	"os/signal"
	"syscall"
)

// LinuxSignalHandler 实现Linux的信号处理
type LinuxSignalHandler struct{}

// NewSignalHandler 创建一个Linux信号处理器
func NewSignalHandler() ShutdownSignalHandler {
	return &LinuxSignalHandler{}
}

// GetShutdownSignals 返回Linux系统上的终止信号
func (h *LinuxSignalHandler) GetShutdownSignals() []os.Signal {
	return []os.Signal{
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // 终止信号
		syscall.SIGHUP,  // 终端关闭
	}
}

// SetupSignalHandler 设置并返回一个信号通道
func (h *LinuxSignalHandler) SetupSignalHandler() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, h.GetShutdownSignals()...)
	return ch
}

// Cleanup 执行清理操作
func (h *LinuxSignalHandler) Cleanup() error {
	return nil
}
