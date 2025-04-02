//go:build windows
// +build windows

package graceful

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/han1eng/go-terminator/pkg/signals"
)

// Start 启动Windows平台的管理器
func (m *manager) Start() error {
	m.startedMu.Lock()
	defer m.startedMu.Unlock()

	if m.started {
		return errors.Errorf("manager already started")
	}

	// 转换为Windows特定的处理器
	winHandler, ok := m.signalHandler.(*signals.WindowsSignalHandler)
	if !ok {
		return errors.Errorf("expected WindowsSignalHandler")
	}

	// 创建命名管道
	err := winHandler.CreateNamedPipe()
	if err != nil {
		return errors.Errorf("failed to create named pipe: %v", err)
	}

	// 设置信号处理
	sigChan := m.signalHandler.SetupSignalHandler()
	m.started = true

	// 启动管道处理
	go winHandler.HandlePipeConnections(m.shutdownChan)

	// 处理信号
	go func() {
		select {
		case sig := <-sigChan:
			fmt.Printf("Received signal: %v\n", sig)
			close(m.shutdownChan)
		case <-m.shutdownChan:
			// 已通过管道或Stop()触发
		}

		// 执行关闭流程
		ctx, cancel := context.WithTimeout(context.Background(), m.opts.GracePeriod)
		defer cancel()

		m.executeHooks(ctx)
		m.signalHandler.Cleanup()
		close(m.doneChan)
	}()

	return nil
}
