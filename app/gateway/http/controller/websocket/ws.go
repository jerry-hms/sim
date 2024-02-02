package websocket

import (
	"github.com/gin-gonic/gin"
	wsSerivce "sim/app/services/websocket"
)

type Ws struct {
}

func (w *Ws) OnOpen(c *gin.Context) (*wsSerivce.Ws, bool) {
	return (&wsSerivce.Ws{}).OnOpen(c)
}

func (w *Ws) OnMessage(ws *wsSerivce.Ws, c *gin.Context) {
	ws.OnMessage()
}

func (w *Ws) Handle(c *gin.Context) {
	if serviceWs, ok := w.OnOpen(c); ok {
		w.OnMessage(serviceWs, c)
	}
}
