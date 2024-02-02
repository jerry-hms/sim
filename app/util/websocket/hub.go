package websocket

import (
	"github.com/google/uuid"
	"sync"
)

var HubIns *Hub
var HubOnce sync.Once

type Hub struct {
	Register   chan *Client
	UnRegister chan *Client
	Clients    map[uuid.UUID]struct {
		Client *Client
		Status bool
	}
}

func CreateHubFactory() *Hub {
	HubOnce.Do(func() {
		HubIns = &Hub{
			Register:   make(chan *Client),
			UnRegister: make(chan *Client),
			Clients: make(map[uuid.UUID]struct {
				Client *Client
				Status bool
			}),
		}
	})

	return HubIns
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.ClientId] = struct {
				Client *Client
				Status bool
			}{Client: client, Status: true}
			break
		case client := <-h.UnRegister:
			c := h.Clients[client.ClientId]
			if c.Status {
				_ = client.Conn.Close()
				delete(h.Clients, client.ClientId)
			}
		}
	}
}

// GetClientByClientId 根据uid获取ws客户端
func (h *Hub) GetClientByClientId(clientId uuid.UUID) *Client {
	client := h.Clients[clientId]
	return client.Client
}
