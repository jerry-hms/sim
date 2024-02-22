package model

import (
	"encoding/json"
	"sim/app/model/common"
	"sim/app/services/im/core"
	"strconv"
)

func CreateImSessionMessageFactory() *ImSessionMessage {
	return &ImSessionMessage{BaseModel: BaseModel{DB: ConnDb()}}
}

type ImSessionMessage struct {
	BaseModel
	RelId       uint64      `json:"rel_id" gorm:"column:rel_id"`
	MessageId   string      `json:"message_id" gorm:"column:message_id"`
	SepSvr      string      `json:"sep_svr" gorm:"column:sep_svr"`
	SenderId    uint64      `json:"sender_id" gorm:"column:sender_id"`
	Sender      common.JSON `json:"sender" gorm:"column:sender;type:json"`
	ReceiverId  uint64      `json:"receiver_id" gorm:"column:receiver_id"`
	SendContent common.JSON `json:"send_content" gorm:"column:send_content;type:json"`
	IsRead      int8        `json:"is_read" gorm:"column:is_read"`
}

// RecordMessage 记录聊天内容
func (m *ImSessionMessage) RecordMessage(rel *ImSessionRelation, sender_id uint64, recv_id uint64, message *core.Message) error {
	sender, _ := json.Marshal(message.Attachments.Sender)
	sendContent, _ := json.Marshal(message.Attachments.Message)
	m.RelId = rel.Id
	m.SenderId = sender_id
	m.MessageId = message.MessageId
	m.Sender = sender
	m.ReceiverId = recv_id
	m.SendContent = sendContent
	// 根据会话的序列号递增1
	m.SepSvr = rel.SepSvr
	m.IncSepSvr()

	result := m.Create(&m)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *ImSessionMessage) IncSepSvr() {
	sepSvr, _ := strconv.Atoi(m.SepSvr)
	m.SepSvr = strconv.Itoa(sepSvr + 1)
}
