package client

import (
	"errors"
	"github.com/google/uuid"
	"sim/app/util/websocket"
	"sync"
)

var Instance *Manage
var InstanceOnce sync.Once

type Manage struct {
	Maps map[uint64]uuid.UUID
}

func CreateManage() *Manage {
	InstanceOnce.Do(func() {
		Instance = &Manage{Maps: make(map[uint64]uuid.UUID)}
	})
	return Instance
}

// BindUidToClientId 绑定uid与client_id
func (m *Manage) BindUidToClientId(uid uint64, clientId uuid.UUID) bool {
	m.Maps[uid] = clientId
	return true
}

func (m *Manage) GetClientId(uid uint64) uuid.UUID {
	return m.Maps[uid]
}

// GetClientIdByUid 根据uid获取ws客户端
func (m *Manage) GetClientIdByUid(uid uint64) (*websocket.Client, error) {
	clientId := m.GetClientId(uid)
	client := websocket.CreateHubFactory().GetClientByClientId(clientId)
	if client == nil {
		return nil, errors.New("客户端不存在")
	}
	return client, nil
}

// Send 发送消息
func (m *Manage) Send(recv_id uint64, message interface{}) error {
	client, err := m.GetClientIdByUid(recv_id)
	if err != nil {
		return err
	}
	return client.SendJsonMessage(message)
}
