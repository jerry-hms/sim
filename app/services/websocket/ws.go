package websocket

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sim/app/global/consts"
	"sim/app/global/variable"
	"sim/app/util/websocket"
)

type Ws struct {
	WsClient *websocket.Client
}

// Send 发送消息
func (c *Ws) Send(recvId uint64, message interface{}) error {
	client, err := websocket.CreateManage().GetClientIdByUid(recvId)
	if err != nil {
		return err
	}
	return client.SendJsonMessage(message)
}

func (w *Ws) OnOpen(context *gin.Context) (*Ws, bool) {
	if client, ok := (&websocket.Client{}).OnOpen(context); ok {
		w.WsClient = client
		// 一旦握手+协议升级成功，就为每一个连接开启一个自动化的隐式心跳检测包
		go w.WsClient.HeartBeat()
		return w, true
	} else {
		return nil, false
	}
}

func (w *Ws) OnMessage() {
	w.WsClient.ReadPump(func(messageType int, message []byte) {
		err := w.WsClient.SendMessage(messageType, string(message))
		if err != nil {
			variable.ZapLog.Error(consts.WebsocketSendMessageFailMsg, zap.Error(err))
		}
	}, w.OnError, w.OnClose)
}

func (w *Ws) OnError(err error) {
	w.WsClient.State = 0 // 发生错误，状态设置为0, 心跳检测协程则自动退出
}

// OnClose 客户端关闭回调，发生onError回调以后会继续回调该函数
func (w *Ws) OnClose() {
	w.WsClient.Hub.UnRegister <- w.WsClient // 向hub管道投递一条注销消息，由hub中心负责关闭连接、删除在线数据
}
