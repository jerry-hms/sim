package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"sim/app/global/variable"
	"sync"
	"time"
)

type Client struct {
	Hub                *Hub            // 负责处理客户端注册、注销、在线管理
	Conn               *websocket.Conn // 一个ws连接
	Send               chan []byte     // 一个ws连接存储自己的消息管道
	PingPeriod         time.Duration
	ReadDeadline       time.Duration
	WriteDeadline      time.Duration
	HeartbeatFailTimes int
	State              uint8     // ws状态，1=ok；0=出错、掉线等
	ClientId           uuid.UUID // 客户端唯一ID
	sync.RWMutex
}

func (c *Client) OnOpen(context *gin.Context) (*Client, bool) {
	defer func() {
		err := recover()
		if err != nil {
			if val, ok := err.(error); ok {
				fmt.Println("websocket报错:", val)
			}
		}
	}()

	var upgrade = websocket.Upgrader{
		ReadBufferSize:  variable.ConfigYml.GetInt("websocket.ReadBufferSize"),
		WriteBufferSize: variable.ConfigYml.GetInt("websocket.WriteBufferSize"),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// 将http协议升级到websocket，然后返回一个长连接客户端
	if wsConn, err := upgrade.Upgrade(context.Writer, context.Request, nil); err != nil {
		fmt.Println("websocket协议升级报错:", err.Error())
		return nil, false
	} else {

		if wsHub, ok := variable.WebsocketHub.(*Hub); ok {
			c.Hub = wsHub
		}
		c.Conn = wsConn

		c.Send = make(chan []byte, variable.ConfigYml.GetInt("websocket.WriteBufferSize"))
		c.PingPeriod = time.Second * variable.ConfigYml.GetDuration("websocket.PingPeriod")
		c.ReadDeadline = time.Second * variable.ConfigYml.GetDuration("websocket.ReadDeadline")
		c.WriteDeadline = time.Second * variable.ConfigYml.GetDuration("websocket.WriteDeadline")
		c.Conn.SetReadLimit(variable.ConfigYml.GetInt64("websocket.MaxMessageSize"))
		c.Hub.Register <- c
		c.State = 1

		c.ClientId = uuid.New()
		if err := c.SendMessage(websocket.TextMessage, fmt.Sprintf(variable.WebsocketHandshakeSuccess, c.ClientId)); err != nil {
			fmt.Println("websocket连接成功，发送打招呼消息失败", err.Error())
		}

		return c, true
	}
}

// 发送消息
func (c *Client) SendMessage(messageType int, message string) error {
	c.Lock()
	defer func() {
		c.Unlock()
	}()
	if err := c.Conn.SetReadDeadline(time.Now().Add(c.WriteDeadline)); err != nil {
		return err
	}

	if err := c.Conn.WriteMessage(messageType, []byte(message)); err != nil {
		return err
	}
	return nil
}

func (c *Client) SendJsonMessage(message interface{}) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return c.SendMessage(websocket.TextMessage, string(msg))
}

// 监听
func (c *Client) ReadPump(callBackOnMessage func(messageType int, message []byte), callBackOnError func(err error), callBackOnClose func()) {
	// 捕获ws的错误，一但有错误则回调给关闭ws的方法
	defer func() {
		err := recover()
		if err != nil {
			if realErr, ok := err.(error); ok {
				fmt.Println("服务器发生错误", realErr)
			}
		}
		callBackOnClose()
	}()

	for {
		if c.State == 1 {
			mt, data, err := c.Conn.ReadMessage()
			if err == nil {
				callBackOnMessage(mt, data)
			} else {
				callBackOnError(err) // 将错误回调给错误处理方法
				break
			}
		} else {
			callBackOnError(errors.New("客户端已下线"))
			break
		}
	}
}

// HeartBeat 心跳包处理方法
func (c *Client) HeartBeat() {
	ticker := time.NewTicker(c.PingPeriod)
	defer func() {
		err := recover()
		if err != nil {
			if val, ok := err.(error); ok {
				fmt.Println("发生错误，监测中断：", val.Error())
			}
		}
		ticker.Stop()
	}()

	if c.ReadDeadline == 0 {
		_ = c.Conn.SetReadDeadline(time.Time{})
	} else {
		_ = c.Conn.SetReadDeadline(time.Now().Add(c.ReadDeadline))
	}
	c.Conn.SetPongHandler(func(receivedPong string) error {
		fmt.Println("收到客户端的pong", receivedPong)
		if c.ReadDeadline > time.Nanosecond {
			_ = c.Conn.SetReadDeadline(time.Now().Add(c.ReadDeadline))
		} else {
			_ = c.Conn.SetReadDeadline(time.Time{})
		}
		return nil
	})

	for {
		select {
		case <-ticker.C:
			if c.State == 1 {
				err := c.SendMessage(websocket.TextMessage, `{type: "Ping"}`)
				if err != nil {
					c.HeartbeatFailTimes++
					if c.HeartbeatFailTimes > variable.ConfigYml.GetInt("websocket.HeartbeatFailMaxTimes") {
						fmt.Println("发送心跳包失败")
						c.State = 0
						c.Hub.UnRegister <- c
						return
					}
				} else {
					// ping通则清空失败次数
					c.HeartbeatFailTimes = 0
				}
			} else {
				return
			}
		}
	}

}
