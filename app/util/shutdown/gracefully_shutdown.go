package shutdown

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func GracefullyShutdown(s *http.Server) {
	// 创建系统信号接收器，接收关闭信号
	done := make(chan os.Signal, 1)
	/**
	os.Interrupt           -> ctrl+c 的信号
	syscall.SIGINT|SIGTERM -> kill 进程时传递给进程的信号
	*/
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	fmt.Println("正在优雅的关闭网关...")
	if err := s.Shutdown(context.Background()); err != nil {
		fmt.Printf("关闭网关报错 %s", err.Error())
	}
}
