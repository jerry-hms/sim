package client

import (
	"sync"
)

var WsClientIns *Client
var WsClientOnce sync.Once

// GetWsClient 获取ws客户端
func GetWsClient() *Client {
	WsClientOnce.Do(func() {
		WsClientIns = &Client{}
	})

	return WsClientIns
}

// Client ws客户端
type Client struct {
}

// Send 发送消息
func (c *Client) Send(recvId uint64, message interface{}) error {
	client, err := CreateManage().GetClientIdByUid(recvId)
	if err != nil {
		return err
	}
	return client.SendJsonMessage(message)
}
