package graceful

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/han1eng/go-terminator/pkg/signals"
	"github.com/pkg/errors"
)

// 通用管理器结构体
type manager struct {
	opts          Options
	hooks         map[string]hookEntry
	hooksMu       sync.RWMutex
	shutdownChan  chan struct{}
	doneChan      chan struct{}
	signalHandler signals.ShutdownSignalHandler
	started       bool
	startedMu     sync.Mutex
}

// newManager 创建一个新的管理器实例
func newManager(opts Options) ShutdownManager {
	m := &manager{
		opts:          opts,
		hooks:         make(map[string]hookEntry),
		shutdownChan:  make(chan struct{}),
		doneChan:      make(chan struct{}),
		signalHandler: signals.NewSignalHandler(),
	}

	return m
}

// RegisterHook 注册一个关闭钩子
func (m *manager) RegisterHook(name string, hook ShutdownHook, timeout time.Duration) error {
	m.hooksMu.Lock()
	defer m.hooksMu.Unlock()

	if _, exists := m.hooks[name]; exists {
		return errors.Errorf("hook with name %s already registered", name)
	}

	m.hooks[name] = hookEntry{
		hook:    hook,
		timeout: timeout,
	}

	return nil
}

// Done 返回一个在关闭完成时关闭的channel
func (m *manager) Done() <-chan struct{} {
	return m.doneChan
}

// Wait 等待关闭完成
func (m *manager) Wait() error {
	<-m.doneChan
	return nil
}

// Stop 手动触发关闭流程
func (m *manager) Stop() error {
	select {
	case <-m.shutdownChan:
		// 已经在关闭中
		return nil
	default:
		close(m.shutdownChan)
		return nil
	}
}

// executeHooks 执行所有注册的钩子
func (m *manager) executeHooks(ctx context.Context) {
	var wg sync.WaitGroup

	m.hooksMu.RLock()
	defer m.hooksMu.RUnlock()

	for name, entry := range m.hooks {
		wg.Add(1)

		go func(name string, hook ShutdownHook, timeout time.Duration) {
			defer wg.Done()

			hookCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			log.Printf("Executing hook: %s\n", name)
			err := hook.Execute(hookCtx)
			if err != nil {
				log.Printf("Hook %s failed: %v\n", name, err)
			} else {
				log.Printf("Hook %s completed successfully\n", name)
			}
		}(name, entry.hook, entry.timeout)
	}

	wg.Wait()
}
