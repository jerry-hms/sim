package websocket

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

var Instance *Manage
var InstanceOnce sync.Once

// Manage websocket客户端管理器
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
func (m *Manage) GetClientIdByUid(uid uint64) (*Client, error) {
	clientId := m.GetClientId(uid)
	client := CreateHubFactory().GetClientByClientId(clientId)
	if client == nil {
		return nil, errors.New("客户端不存在")
	}
	return client, nil
}
