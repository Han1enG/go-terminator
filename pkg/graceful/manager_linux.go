//go:build linux
// +build linux

package graceful

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

// Start 启动Linux平台的管理器
func (m *manager) Start() error {
	m.startedMu.Lock()
	defer m.startedMu.Unlock()

	if m.started {
		return errors.Errorf("manager already started")
	}

	// 设置信号处理
	sigChan := m.signalHandler.SetupSignalHandler()
	m.started = true

	go func() {
		select {
		case sig := <-sigChan:
			fmt.Printf("Received signal: %v\n", sig)
			close(m.shutdownChan)
		case <-m.shutdownChan:
			// 已通过Stop()手动触发
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
