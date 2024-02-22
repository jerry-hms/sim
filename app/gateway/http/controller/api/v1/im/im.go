package im

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sim/app/gateway/http/controller/api"
	"sim/app/gateway/rpc"
	"sim/app/global/variable"
	"sim/app/util/jwt"
	"sim/app/util/response"
	"sim/app/util/websocket"
	pb "sim/idl/im"
)

type ImControl struct {
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
	manage, _ := variable.WebsocketManage.(*websocket.Manage)
	if ok := manage.BindUidToClientId(u.Info.Id, cid); !ok {
		response.Fail(ctx, "绑定失败", nil)
		return
	}

	response.Success(ctx, "绑定成功", nil)
}

// Send 聊天消息发送
func (c *ImControl) Send(ctx *gin.Context) {
	user := api.GetLoginUser(ctx)
	req := &pb.SendRequest{
		SenderId: user.Id,
		RecvId:   uint64(ctx.GetFloat64("recv_id")),
		Content:  ctx.GetString("content"),
		Type:     ctx.GetString("type"),
		Scene:    ctx.GetString("scene"),
		Url:      ctx.GetString("url"),
		Width:    int64(ctx.GetFloat64("width")),
		Height:   int64(ctx.GetFloat64("height")),
	}

	send, err := rpc.ImClient.Send(ctx, req)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, "发送成功", send)
}

func (c *ImControl) SessionList(ctx *gin.Context) {
	user := api.GetLoginUser(ctx)
	req := &pb.SessionListRequest{
		UserId:   user.Id,
		Page:     int64(ctx.GetFloat64("page")),
		PageSize: int64(ctx.GetFloat64("page_size")),
	}
	resp, _ := rpc.ImClient.SessionList(ctx, req)
	response.Success(ctx, "操作成功", resp)
}
