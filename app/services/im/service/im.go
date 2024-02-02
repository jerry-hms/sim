package service

import (
	"context"
	"sim/app/services/im/rpc"
	"sim/app/services/websocket/client"
	pb "sim/idl/pb/im"
	userPb "sim/idl/pb/user"
	"sync"
)

var ImSrv *Im
var ImIns sync.Once

func GetImSrv() *Im {
	ImIns.Do(func() {
		ImSrv = &Im{}
	})

	return ImSrv
}

type Im struct {
	*pb.UnimplementedImServiceServer
	MessageId string
}

// Send 发送消息
func (i *Im) Send(c context.Context, req *pb.SendRequest) (*pb.SendResponse, error) {
	userReq := &userPb.UserInfoRequest{
		Id: req.RecvId,
	}
	user, err := rpc.UserClient.UserInfo(c, userReq)
	if err != nil {
		return nil, err
	}
	wsClientManage := client.CreateManage()
	err = wsClientManage.Send(req.RecvId, req.String())
	if err != nil {
		return nil, err
	}
	return nil, nil
}
