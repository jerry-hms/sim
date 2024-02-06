package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sim/app/gateway/rpc"
	"sim/app/global/variable"
	"sim/app/util/rabbitmq"
	"sim/app/util/shutdown"
	_ "sim/bootstrap"
	"sim/routes"
	"syscall"
	"time"
)

func main() {
	// 初始化rpc服务
	rpc.Init()
	go runServer()
	fmt.Printf("gateway listen on %s\n", variable.ConfigYml.GetString("server.port"))

	// 启动mq
	go rabbitmq.InitRabbitMq()
	// 优雅的关闭服务
	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGKILL)
		s := <-osSignals
		fmt.Println("exit! ", s)
	}
}

func runServer() {
	r := routes.RegisterRoute()
	server := &http.Server{
		Addr:           variable.ConfigYml.GetString("server.port"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("网关启动报错: %s\n", err.Error())
	}

	go func() {
		shutdown.GracefullyShutdown(server)
	}()
}
