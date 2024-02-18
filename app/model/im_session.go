package model

import (
	"encoding/json"
	"sim/app/global/variable"
	"sim/app/model/common"
	"sim/app/services/im/core"
	"sim/app/services/im/core/interf"
)

func CreateImSessionFactory() *ImSession {
	return &ImSession{BaseModel: BaseModel{DB: ConnDb()}}
}

// 从消息体中解析出最后发送的消息内容
func parseLastMessage(message *core.Message) string {
	msgType, exists := message.Attachments.Message["type"]
	if !exists {
		variable.ZapLog.Error("message type field is not defined")
	}
	mt, _ := msgType.(string)
	x, _ := message.Attachments.Message[mt]
	ins, ok := x.(interf.TypeInterface)
	if !ok {
		return ""
	}
	return ins.ParseContent()
}

type ImSession struct {
	BaseModel
	UserId              uint64      `gorm:"column:user_id" json:"user_id"`
	RelId               uint64      `gorm:"column:rel_id" json:"rel_id"`
	SessionName         string      `gorm:"column:session_name" json:"session_name"`
	SepSvr              string      `gorm:"column:sep_svr" json:"sep_svr"`
	LastSenderInfo      common.JSON `gorm:"column:last_sender_info;type:json" json:"last_sender_info"`
	LastMessage         string      `gorm:"column:last_message" json:"last_message"`
	UnreadMessageNumber int16       `gorm:"column:unread_message_number" json:"unread_message_number"`
}

// GetSessionOrCreate 获取会话，如果会话不存在就创建，并且更新会话数据
func (i *ImSession) GetSessionOrCreate(user_id uint64, rel_id uint64, name string, message *core.Message, unread_message_number int16) error {
	senderInfo, _ := json.Marshal(message.Attachments.Sender)

	i.UserId = user_id
	i.RelId = rel_id
	i.SessionName = name
	i.LastSenderInfo = senderInfo
	i.LastMessage = parseLastMessage(message)
	i.UnreadMessageNumber = unread_message_number

	result := i.Where("user_id = ?", user_id).Where("rel_id", rel_id).FirstOrCreate(&i)
	if result.Error != nil {
		return result.Error
	}

	i.SessionName = name
	i.LastSenderInfo = senderInfo
	i.LastMessage = parseLastMessage(message)
	i.UnreadMessageNumber = i.UnreadMessageNumber + unread_message_number

	result = i.Save(&i)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
