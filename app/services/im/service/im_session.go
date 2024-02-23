package service

import (
	"encoding/json"
	"fmt"
	"sim/app/global/consts"
	"sim/app/global/variable"
	"sim/app/model/common"
	"sim/app/services/im/core"
	"sim/app/services/im/core/interf"
	"sim/app/util/pagination"
	"time"
)

func CreateSessionServiceFactory() *Session {
	return &Session{}
}

func getHashDb() {

}

func GetSessionCacheKey(user_id uint64) string {
	return fmt.Sprintf(consts.ImSessionKeyFormat, user_id)
}

func GetSessionSortKey(user_id uint64) string {
	return fmt.Sprintf(consts.ImSessionTopicSortKeyFormat, user_id)
}

// GetSessionCacheField 获取会话redis缓存field字段
func GetSessionCacheField(rel_id uint64) string {
	return fmt.Sprintf(consts.ImSessionFieldFormat, rel_id)
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

type Session struct {
	UserId              uint64      `redis:"user_id" `
	RelId               uint64      `redis:"rel_id" json:"rel_id"`
	SessionName         string      `redis:"session_name" json:"session_name"`
	SepId               int64       `redis:"sep_id" json:"sep_id"`
	LastSenderInfo      common.JSON `redis:"last_sender_info;type:json" json:"last_sender_info"`
	LastMessage         string      `redis:"last_message" json:"last_message"`
	UnreadMessageNumber int16       `redis:"unread_message_number" json:"unread_message_number"`
	CreatedAt           time.Time   `redis:"created_at" json:"created_at"`
	UpdatedAt           time.Time   `redis:"updated_at" json:"updated_at"`
}

func (s *Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

// GetSessionOrCreate 获取会话，如果会话不存在就创建，并且更新会话数据
func (s *Session) GetSessionOrCreate(user_id uint64, rel_id uint64, name string, message *core.Message, unread_message_number int16) error {
	rdb := pagination.GetHashDb().SetHashTable(GetSessionCacheKey(user_id)).SetHashField(GetSessionCacheField(rel_id))

	result, err := rdb.First()
	if err != nil {
		return err
	}

	senderInfo, _ := json.Marshal(message.Attachments.Sender)
	s.UserId = user_id
	s.RelId = rel_id
	s.SessionName = name
	s.LastSenderInfo = senderInfo
	s.LastMessage = parseLastMessage(message)
	s.UpdatedAt = time.Now()
	if result != "" {
		_ = json.Unmarshal([]byte(result), s)
		s.UnreadMessageNumber = s.UnreadMessageNumber + unread_message_number
		s.UpdatedAt = time.Now()
	} else {
		s.UnreadMessageNumber = unread_message_number
		s.CreatedAt = time.Now()
	}
	err = rdb.SetHashTable(GetSessionCacheKey(user_id)).SetHashField(GetSessionCacheField(rel_id)).Store(s)
	return err
}

// List 获取会话列表数据
func (s *Session) List(user_id uint64, page int, page_size int, sort string) (*pagination.Pagination, error) {
	var sessions []*Session
	var session *Session
	paginate := pagination.Pagination{}

	rp, err := pagination.GetHashDb().SetHashTable(GetSessionCacheKey(user_id)).Paginate(page, page_size, "desc")
	if err != nil {
		return nil, err
	}
	for _, item := range rp.Rows {
		_ = json.Unmarshal([]byte(item), &session)
		sessions = append(sessions, session)
	}
	paginate.Rows = sessions
	paginate.Page = rp.Page
	paginate.PageSize = rp.Limit
	paginate.TotalRows = rp.TotalRows
	paginate.TotalPages = rp.TotalPages
	return &paginate, nil
}
