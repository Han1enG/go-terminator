//go:build windows
// +build windows

package signals

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Microsoft/go-winio"
	"github.com/pkg/errors"
)

const (
	// PipeName 是Windows命名管道的名称
	PipeName = `\\.\pipe\graceful-shutdown`
)

// WindowsSignalHandler 实现Windows的信号处理
type WindowsSignalHandler struct {
	pipeListener net.Listener
}

// NewSignalHandler 创建一个Windows信号处理器
func NewSignalHandler() ShutdownSignalHandler {
	return &WindowsSignalHandler{}
}

// GetShutdownSignals 返回Windows系统上的终止信号
func (h *WindowsSignalHandler) GetShutdownSignals() []os.Signal {
	return []os.Signal{
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // 终止信号
	}
}

// CreateNamedPipe 创建Windows命名管道
func (h *WindowsSignalHandler) CreateNamedPipe() error {
	pipeCfg := &winio.PipeConfig{
		SecurityDescriptor: "D:P(A;;GA;;;BA)(A;;GA;;;SY)",
		MessageMode:        true,
		InputBufferSize:    1024,
		OutputBufferSize:   1024,
	}

	var err error
	h.pipeListener, err = winio.ListenPipe(PipeName, pipeCfg)
	return err
}

// SetupSignalHandler 设置并返回一个信号通道
func (h *WindowsSignalHandler) SetupSignalHandler() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, h.GetShutdownSignals()...)
	return ch
}

// HandlePipeConnections 处理来自命名管道的连接
func (h *WindowsSignalHandler) HandlePipeConnections(shutdownChan chan struct{}) {
	if h.pipeListener == nil {
		fmt.Println("Pipe listener not initialized")
		return
	}

	for {
		conn, err := h.pipeListener.Accept()
		if err != nil {
			// 如果已经关闭监听器，就退出循环
			select {
			case <-shutdownChan:
				return
			default:
				fmt.Printf("Pipe accept error: %v\n", err)
				continue
			}
		}

		go func(conn io.ReadWriteCloser) {
			defer conn.Close()

			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("Pipe read error: %v\n", err)
				return
			}

			if string(buf[:n]) == "SHUTDOWN" {
				fmt.Println("Received shutdown command through pipe")
				close(shutdownChan)
			}
		}(conn)
	}
}

// SendShutdownSignal 通过命名管道发送终止信号
func SendShutdownSignal() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := winio.DialPipeContext(ctx, PipeName)
	if err != nil {
		return errors.Errorf("failed to connect to pipe: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("SHUTDOWN"))
	if err != nil {
		return errors.Errorf("failed to send shutdown signal: %v", err)
	}

	return nil
}

// Cleanup 执行清理操作
func (h *WindowsSignalHandler) Cleanup() error {
	if h.pipeListener != nil {
		return h.pipeListener.Close()
	}
	return nil
}
