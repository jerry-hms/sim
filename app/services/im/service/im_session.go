package service

import (
	"context"
	"encoding/json"
	"fmt"
	baseRedis "github.com/go-redis/redis/v8"
	"sim/app/global/consts"
	"sim/app/global/variable"
	"sim/app/model/common"
	"sim/app/services/im/core"
	"sim/app/services/im/core/interf"
	"sim/app/util/pagination"
	"sim/app/util/redis"
	"sim/app/util/tools"
	"time"
)

func CreateSessionServiceFactory() *Session {
	return &Session{}
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
	senderInfo, _ := json.Marshal(message.Attachments.Sender)
	rdb := redis.ConnRedis()
	ets := rdb.HExists(context.Background(), GetSessionCacheKey(user_id), GetSessionCacheField(rel_id))
	var err error
	if !ets.Val() {
		s.UserId = user_id
		s.RelId = rel_id
		s.SessionName = name
		s.LastSenderInfo = senderInfo
		s.LastMessage = parseLastMessage(message)
		s.UnreadMessageNumber = unread_message_number
		s.CreatedAt = time.Now()
		s.UpdatedAt = time.Now()
		jsonStr, _ := json.Marshal(s)

		// 数据存入redis
		if err = rdb.HSet(context.Background(), GetSessionCacheKey(user_id), GetSessionCacheField(rel_id), string(jsonStr)).Err(); err != nil {
			return err
		}
		// 再存入对应的排序方便做分页查询
		sort := rdb.HLen(context.Background(), GetSessionCacheKey(user_id)).Val()
		item := &baseRedis.Z{
			Score:  float64(sort),
			Member: GetSessionCacheField(rel_id),
		}
		rdb.ZAdd(context.Background(), GetSessionSortKey(user_id), item)
	} else {
		if err = rdb.HGet(context.Background(), GetSessionCacheKey(user_id), GetSessionCacheField(rel_id)).Scan(s); err != nil {
			return err
		}
		s.SessionName = name
		s.LastSenderInfo = senderInfo
		s.LastMessage = parseLastMessage(message)
		s.UnreadMessageNumber = s.UnreadMessageNumber + unread_message_number
		s.UpdatedAt = time.Now()
		jsonStr, _ := json.Marshal(s)

		if err = rdb.HSet(context.Background(), GetSessionCacheKey(user_id), GetSessionCacheField(rel_id), string(jsonStr)).Err(); err != nil {
			return err
		}
	}

	return nil
}

// List 获取会话列表数据
func (s *Session) List(user_id uint64, p interface{}) (*pagination.Pagination, error) {
	var sessions []*Session
	var session Session
	var paginate pagination.Pagination

	_ = tools.InterfaceToStruct(p, paginate)
	rdb := redis.ConnRedis()
	// 计算分页总页数
	paginate.CountTotalPages(rdb.HLen(context.Background(), GetSessionCacheKey(user_id)).Val())
	// 获取会话的排序,并且分页
	keys, _ := rdb.ZRevRange(context.Background(), GetSessionSortKey(user_id), int64(paginate.GetOffset()), int64(paginate.GetEnd())).Result()
	for _, key := range keys {
		// 从redis中取出会话数据
		if err := rdb.HGet(context.Background(), GetSessionCacheKey(user_id), key).Scan(&session); err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}
	paginate.Rows = sessions
	return &paginate, nil
}
