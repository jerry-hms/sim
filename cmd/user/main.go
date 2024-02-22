package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sim/app/global/variable"
	userSrv "sim/app/services/user/service"
	"sim/app/util/discovery"
	_ "sim/bootstrap"
	pb "sim/idl/user"
)

var Port = "10001"

func main() {
	flag.Parse()
	// etcd服务地址
	etcdAddress := []string{variable.ConfigYml.GetString("etcd.address")}
	// 服务注册
	etcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	defer etcdRegister.Stop()

	// user rpc服务开启
	grpcAddress := fmt.Sprintf("%s:%s",
		variable.ConfigYml.GetStringSlice("services.user.host"),
		Port)
	userNode := discovery.Server{
		Name: variable.ConfigYml.GetString("services.user.name"),
		Addr: grpcAddress,
	}
	s := grpc.NewServer()
	defer s.Stop()

	pb.RegisterUserServiceServer(s, userSrv.GetUserSrv(grpcAddress))
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	// 将当前user rpc服务注册到etcd中
	if _, err = etcdRegister.Register(userNode, 10); err != nil {
		panic(fmt.Sprintf("start server failed, err: %v", err))
	}
	if err = s.Serve(lis); err != nil {
		panic(err)
	}

	fmt.Println("启动成功")
}
