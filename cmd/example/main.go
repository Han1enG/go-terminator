package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/han1eng/go-terminator/pkg/graceful"
)

func main() {
	// 创建HTTP服务器
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello, World!")
		}),
	}

	// 创建终止管理器
	opts := graceful.DefaultOptions()
	manager := graceful.NewManager(opts)

	// 注册HTTP服务器关闭钩子
	err := manager.RegisterHook("http-server", graceful.HookFunc(func(ctx context.Context) error {
		log.Println("Shutting down HTTP server...")
		return server.Shutdown(ctx)
	}), 10*time.Second)

	if err != nil {
		log.Fatalf("Failed to register hook: %v\n", err)
	}

	// 启动终止管理器
	err = manager.Start()
	if err != nil {
		log.Fatalf("Failed to start manager: %v\n", err)
	}

	log.Println("Starting HTTP server on :8080")

	// 在单独的goroutine中启动HTTP服务器
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v\n", err)
		}
	}()

	log.Println("Server is running. Press Ctrl+C to stop.")

	// 等待终止信号
	manager.Wait()

	log.Println("Graceful shutdown completed")
}
