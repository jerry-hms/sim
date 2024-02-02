package rpc

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"sim/app/global/variable"
	"sim/app/util/discovery"
	"sim/idl/pb/user"
	"time"
)

var (
	Register   *discovery.Resolver
	ctx        context.Context
	CancelFunc context.CancelFunc

	UserClient user.UserServiceClient // user rpc服务
)

// Init 初始化rpc
func Init() {
	Register = discovery.NewResolver([]string{variable.ConfigYml.GetString("etcd.address")}, logrus.New())
	resolver.Register(Register)
	ctx, CancelFunc = context.WithTimeout(context.Background(), 3*time.Second)

	defer Register.Close()
	InitClient(variable.ConfigYml.GetString("services.user.name"), &UserClient)
}

// InitClient 初始化客户端
func InitClient(serviceName string, client interface{}) {
	conn, err := connectServer(serviceName)
	if err != nil {
		panic(err)
	}

	switch c := client.(type) {
	case *user.UserServiceClient:
		*c = user.NewUserServiceClient(conn)
	default:
		panic("unsupported client type")
	}
}

// 连接服务
func connectServer(serviceName string) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	addr := fmt.Sprintf("%s:///%s", Register.Scheme(), serviceName)

	if variable.ConfigYml.GetBool(fmt.Sprintf("services.%s.loadBalance", serviceName)) {
		log.Printf("load balance enabled for %s\n", serviceName)
		opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	}

	conn, err = grpc.DialContext(ctx, addr, opts...)
	return
}
