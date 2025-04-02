package graceful

import (
	"time"
)

// Options 配置选项
type Options struct {
	// GracePeriod 是执行钩子函数的总超时时间
	GracePeriod time.Duration

	// ForceTimeout 是强制终止前的等待时间
	ForceTimeout time.Duration
}

// DefaultOptions 返回默认配置选项
func DefaultOptions() Options {
	return Options{
		GracePeriod:  30 * time.Second,
		ForceTimeout: 5 * time.Second,
	}
}
