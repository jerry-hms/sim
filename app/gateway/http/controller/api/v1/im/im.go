package im

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sim/app/gateway/rpc"
	"sim/app/services/websocket/client"
	"sim/app/util/jwt"
	"sim/app/util/response"
	pb "sim/idl/pb/im"
)

type ImControl struct {
}

func (i *ImControl) Ws() {
	if serviceWs, ok := w.OnOpen(c); ok {
		w.OnMessage(serviceWs, c)
	}
}

func (c *ImControl) BindToWs(ctx *gin.Context) {
	clientId := ctx.GetString("client_id")
	jwtUser, _ := ctx.Get("user")

	u, _ := jwtUser.(*jwt.Claims)
	cid, err := uuid.Parse(clientId)
	if err != nil {
		response.Fail(ctx, "client_id不合法", nil)
		return
	}
	wsClientManage := client.CreateManage()
	if ok := wsClientManage.BindUidToClientId(u.Info.Id, cid); !ok {
		response.Fail(ctx, "绑定失败", nil)
		return
	}

	response.Success(ctx, "绑定成功", nil)
}

// Send 聊天消息发送
func (c *ImControl) Send(ctx *gin.Context) {
	req := &pb.SendRequest{
		RecvId:  uint64(ctx.GetFloat64("recv_id")),
		Content: ctx.GetString("content"),
		Type:    ctx.GetString("type"),
		Scene:   ctx.GetString("scene"),
		Url:     ctx.GetString("url"),
		Width:   int64(ctx.GetFloat64("width")),
		Height:  int64(ctx.GetFloat64("height")),
	}

	send, err := rpc.ImClient.Send(ctx, req)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, "发送成功", send)
}
