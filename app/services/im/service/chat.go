package service

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sim/app/global/variable"
	"sim/app/model"
	imCore "sim/app/services/im/core"
	"sim/app/services/im/rpc"
	userPb "sim/idl/user"
)

func CreateChat() *Chat {
	return &Chat{}
}

type Chat struct {
}

// Say 发消息
func (c *Chat) Say(user *userPb.UserResponse, recv_id uint64, message *imCore.Message) {
	db := model.ConnDb()
	err := db.Transaction(func(tx *gorm.DB) error {
		// 建立会话关系
		rel := model.CreateImSessionRelationFactory().GetRelationOrCreate(user.GetId(), recv_id, message.Scene)
		// 处理不同场景的聊天逻辑
		var err error
		switch message.Scene {
		case "friend":
			err = c.toFriend(user, recv_id, rel, message)
		case "group":
			// toGroup...
		}
		if err != nil {
			return err
		}
		// 会话的总序列号增加
		rel.IncSepSvr()
		return nil
	})

	if err != nil {
		variable.ZapLog.Error("会话逻辑处理错误", zap.Error(err))
	}
}

// 处理好友聊天逻辑
func (c *Chat) toFriend(user *userPb.UserResponse, recv_id uint64, rel *model.ImSessionRelation, message *imCore.Message) error {
	// 获取接收者信息
	recv, err := rpc.UserClient.UserInfo(context.Background(), &userPb.UserInfoRequest{
		Id: recv_id,
	})
	if err != nil {
		return err
	}
	// 创建发送人会话
	err = CreateSessionServiceFactory().GetSessionOrCreate(user.GetId(), rel.Id, recv.NickName, message, 0)
	if err != nil {
		return errors.New("创建发送者会话报错:" + err.Error())
	}
	// 创建接收人会话
	err = CreateSessionServiceFactory().GetSessionOrCreate(recv_id, rel.Id, user.GetNickName(), message, 1)
	if err != nil {
		return errors.New("创建接收者会话报错:" + err.Error())
	}
	// 记录发送消息
	err = model.CreateImSessionMessageFactory().RecordMessage(rel, user.Id, recv_id, message)
	if err != nil {
		return errors.New("记录发送消息报错:" + err.Error())
	}

	return nil
}
