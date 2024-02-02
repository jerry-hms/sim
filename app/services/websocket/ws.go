package websocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sim/app/util/websocket"
)

type Ws struct {
	WsClient *websocket.Client
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
			fmt.Println("消息发送出现错误", err.Error())
		}
	}, w.OnError, w.OnClose)
}

func (w *Ws) OnError(err error) {
	w.WsClient.State = 0 // 发生错误，状态设置为0, 心跳检测协程则自动退出
	//fmt.Printf("远端掉线、卡死、刷新浏览器等会触发该错误: %v\n", err.Error())
}

// OnClose 客户端关闭回调，发生onError回调以后会继续回调该函数
func (w *Ws) OnClose() {
	w.WsClient.Hub.UnRegister <- w.WsClient // 向hub管道投递一条注销消息，由hub中心负责关闭连接、删除在线数据
}