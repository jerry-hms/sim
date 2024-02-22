package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sim/app/global/variable"
	imCore "sim/app/services/im/core"
	"sim/app/services/im/rpc"
	"sim/app/util/queue"
	"sim/idl/im"
	userPb "sim/idl/user"
	"sync"
)

var ImSrv *Im
var ImIns sync.Once

func GetImSrv(port string) *Im {
	ImIns.Do(func() {
		ImSrv = &Im{
			Port: port,
		}
	})

	return ImSrv
}

// Im 即时聊天服务
type Im struct {
	*im.UnimplementedImServiceServer
	Port string
}

// Send 发送消息
func (i *Im) Send(c context.Context, req *im.SendRequest) (*im.SendResponse, error) {
	// 获取发送人信息
	user, err := rpc.UserClient.UserInfo(c, &userPb.UserInfoRequest{
		Id: req.SenderId,
	})
	if err != nil {
		return nil, err
	}
	// 声明发送人
	sender := struct {
		UserId   uint64 `json:"user_id"`
		NickName string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}{
		UserId:   user.Id,
		NickName: user.NickName,
		Avatar:   user.Avatar,
	}
	// 组合消息结构
	message := imCore.CreateMessage().SetScene(req.Scene).SetSender(sender).SetMessage(req).Parse()
	str, _ := json.Marshal(message)
	qMsg := fmt.Sprintf(`{"recv_id": %d, "message": %s}`, req.RecvId, string(str))
	// 将消息投递到队列中
	s := queue.Pusher("send_message", qMsg, "publish_subscribe", 0)
	if s {
		variable.ZapLog.Info(fmt.Sprintf("服务[%s]消息投递成功", i.Port))
	}
	// 处理会话数据
	go CreateChat().Say(user, req.RecvId, message)

	sendResp := &im.SendResponse{}
	sendResp.MessageId = message.MessageId

	return sendResp, nil
}

// SessionList 会话列表数据
func (i *Im) SessionList(ctx context.Context, req *im.SessionListRequest) (*im.SessionListResponse, error) {
	list, err := CreateSessionServiceFactory().List(req.UserId, req)
	if err != nil {
		return nil, err
	}
	response := &im.SessionListResponse{}
	response.Page = int64(list.Page)
	response.PageSize = int64(list.PageSize)
	response.TotalRows = list.TotalRows
	response.TotalPages = int64(list.TotalPages)
	rows, _ := list.Rows.([]*Session)
	for _, item := range rows {
		var info *im.LastSenderInfo
		_ = json.Unmarshal(item.LastSenderInfo, &info)
		response.Rows = append(response.Rows, &im.Session{
			RelId:               item.RelId,
			SessionName:         item.SessionName,
			SepId:               item.SepId,
			LastSenderInfo:      info,
			LastMessage:         item.LastMessage,
			UnreadMessageNumber: uint32(item.UnreadMessageNumber),
			CreatedAt:           item.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return response, nil
}
