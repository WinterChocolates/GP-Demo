package cmd

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func StartServer(router *gin.Engine) {
	port := viper.GetString("server.port")
	address := ":" + port

	server := &http.Server{
		Addr:    address,
		Handler: router,
	}

	// 创建错误通道和信号通道
	serverErr := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 启动服务器协程
	go func() {
		log.Printf("🚀 服务器启动中，监听地址: %s", address)
		err := server.ListenAndServe()
		serverErr <- err // 直接传递错误，不做提前判断
	}()

	// 等待信号或错误
	select {
	case err := <-serverErr:
		if err != nil {
			// 使用 errors.Is 检查错误是否为 http.ErrServerClosed
			if errors.Is(err, http.ErrServerClosed) {
				log.Println("✅ 服务器正常关闭")
			} else {
				log.Fatalf("❌ 服务器启动失败: %v", err)
			}
		} else {
			log.Println("✅ 服务器正常关闭")
		}
	case sig := <-sigChan:
		log.Printf("🛑 收到信号 %s，开始优雅关闭...", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("❌ 服务器关闭异常: %v", err)
		} else {
			log.Println("✅ 服务器优雅关闭完成")
		}
	}
}
