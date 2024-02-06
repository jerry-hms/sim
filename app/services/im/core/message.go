package core

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"sim/app/services/im/core/types"
	"sync"
	"time"
)

func CreateMessage() *Message {
	return &Message{}
}

// Message 推送消息体
type Message struct {
	Scene       string       `json:"scene"`
	Attachments *Attachments `json:"attachments"`
	Time        int64        `json:"time"`
	MessageId   string       `json:"message_id"`
}

// SetScene 设置发送场景
func (m *Message) SetScene(scene string) *Message {
	m.Scene = scene
	return m
}

// 设置Attachments
func (m *Message) setAttachments() {
	m.Attachments = createAttachments()
}

// 设置发送时间
func (m *Message) setTime() {
	m.Time = time.Now().Unix()
}

// 生成消息ID
func (m *Message) generateMessageId() {
	m.MessageId = uuid.NewString()
}

// GetMessageId 获取消息ID
func (m *Message) GetMessageId() string {
	return m.MessageId
}

// SetSender 设置sender数据
func (m *Message) SetSender(sender interface{}) *Message {
	if m.Attachments == nil {
		m.setAttachments()
	}
	m.Attachments.setSender(sender)
	return m
}

// SetMessage 设置消息内容
func (m *Message) SetMessage(message interface{}) *Message {
	if m.Attachments == nil {
		m.setAttachments()
	}
	var msg map[string]interface{}
	str, _ := json.Marshal(message)
	_ = json.Unmarshal(str, &msg)
	_ = m.Attachments.parseMessage(msg)
	return m
}

// SetExtra 设置其他参数
func (m *Message) SetExtra(data interface{}) *Message {
	if values, ok := data.(map[string]interface{}); ok {
		if m.Attachments == nil {
			m.setAttachments()
		}
		for k, v := range values {
			m.Attachments.Extra[k] = v
		}
	}
	return m
}

// Parse 解析消息体
func (m *Message) Parse() *Message {
	m.setAttachments()
	m.setTime()
	m.generateMessageId()

	return m
}

var attIns *Attachments
var attOnce sync.Once

func createAttachments() *Attachments {
	attOnce.Do(func() {
		attIns = &Attachments{
			Sender:  make(map[string]interface{}),
			Message: make(map[string]interface{}),
			Extra:   make(map[string]interface{}),
		}
	})
	return attIns
}

type Attachments struct {
	Sender  map[string]interface{} `json:"sender"`
	Message map[string]interface{} `json:"message"`
	Extra   map[string]interface{} `json:"extra"`
}

// 设置sender
func (a *Attachments) setSender(sender interface{}) *Attachments {
	str, _ := json.Marshal(sender)
	_ = json.Unmarshal(str, &a.Sender)
	return a
}

// 解析message，并且填充给Attachments
func (a *Attachments) parseMessage(message map[string]interface{}) error {
	t, ets := message["type"]
	if !ets {
		return errors.New("message type is empty")
	}
	ty, _ := t.(string)
	x, ets := types.MessageTypes[ty]
	if !ets {
		return errors.New("message type is not defined")
	}
	a.Message["type"] = ty
	strut, err := x.ParseParams(message)
	if err != nil {
		return err
	}
	a.Message[ty] = strut
	return nil
}

// 设置其他参数
func (a *Attachments) setExtra(key string, value interface{}) {
	a.Message[key] = value
}
