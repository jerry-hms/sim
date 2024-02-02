package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sim/app/gateway/rpc"
	"sim/app/services/websocket/client"
	"sim/app/util/jwt"
	"sim/app/util/response"
	pb "sim/idl/pb/im"
)

type Chat struct {
}

func (c *Chat) BindToWs(ctx *gin.Context) {
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

	wsClientManage.Send(u.Info.Id, gin.H{
		"msg": "成功建立绑定",
	})
	if err != nil {
		response.Fail(ctx, "发送失败", nil)
		return
	}
	response.Success(ctx, "发送成功", nil)
}

func (c *Chat) Send(ctx *gin.Context) {
	req := &pb.SendRequest{
		RecvId:  uint64(ctx.GetFloat64("recv_id")),
		Content: ctx.GetString("content"),
		Type:    ctx.GetString("type"),
		Scene:   ctx.GetString("scene"),
	}

	send, err := rpc.ImClient.Send(ctx, req)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, "发送成功", send)
}
