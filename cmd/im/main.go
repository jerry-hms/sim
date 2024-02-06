package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sim/app/global/variable"
	"sim/app/services/im/rpc"
	imSrv "sim/app/services/im/service"
	"sim/app/util/discovery"
	_ "sim/bootstrap"
	"sim/idl/pb/im"
)

var (
	port = flag.String("port", "20002", "listening port")
)

func main() {
	rpc.Init()
	flag.Parse()
	// etcd服务地址
	etcdAddress := []string{variable.ConfigYml.GetString("etcd.address")}
	// 服务注册
	etcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	defer etcdRegister.Stop()

	grpcAddress := fmt.Sprintf("%s:%s",
		variable.ConfigYml.GetStringSlice("services.im.host"),
		*port)
	chatNode := discovery.Server{
		Name: variable.ConfigYml.GetString("services.im.name"),
		Addr: grpcAddress,
	}
	s := grpc.NewServer()
	defer s.Stop()

	// 将im rpc服务注册到对应的service中
	im.RegisterImServiceServer(s, imSrv.GetImSrv())
	lis, err := net.Listen("tcp", grpcAddress)

	if err != nil {
		panic(err)
	}
	// 将当前user rpc服务注册到etcd中
	if _, err := etcdRegister.Register(chatNode, 10); err != nil {
		panic(fmt.Sprintf("start server failed, err: %v", err))
	}
	if err := s.Serve(lis); err != nil {
		panic(err)
	}

}
