package service

import (
	"context"
	"encoding/json"
	"fmt"
	imCore "sim/app/services/im/core"
	"sim/app/services/im/rpc"
	"sim/app/util/queue"
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

// Im 即时聊天服务
type Im struct {
	*pb.UnimplementedImServiceServer
}

// Send 发送消息
func (i *Im) Send(c context.Context, req *pb.SendRequest) (*pb.SendResponse, error) {
	userReq := &userPb.UserInfoRequest{
		Id: req.RecvId,
	}
	// 获取发送人信息
	user, err := rpc.UserClient.UserInfo(c, userReq)
	if err != nil {
		return nil, err
	}

	// 组合消息结构
	message := imCore.CreateMessage().SetScene(req.Scene).SetSender(user).SetMessage(req).Parse()
	str, _ := json.Marshal(message)
	qMsg := fmt.Sprintf(`{"recv_id": %d, "message": %s}`, req.RecvId, string(str))
	// 将消息投递到队列中
	s := queue.Pusher("send_message", qMsg, "publish_subscribe", 0)
	if s {
		fmt.Println("消息投递成功")
	}
	sendResp := &pb.SendResponse{}
	sendResp.MessageId = message.MessageId

	return sendResp, nil
}
