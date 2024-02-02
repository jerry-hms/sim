package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sim/app/util/encrypt"
)

type UserModel struct {
	BaseModel
	Username string `gorm:"column:username;size:20" json:"username"`
	Password string `gorm:"column:password;size:128" json:"password"`
	Nickname string `gorm:"column:nickname;size:20" json:"nickname"`
	Avatar   string `gorm:"column:avatar" json:"avatar"`
	Mobile   string `gorm:"column:mobile;size:11" json:"mobile"`
}

func CreateUserFactory() *UserModel {
	return &UserModel{BaseModel: BaseModel{DB: ConnDb()}}
}

func (u *UserModel) TableName() string {
	return "user"
}

// 注册
func (u *UserModel) Register(user *UserModel) error {
	var curUser UserModel

	u.Where("username = ?", user.Username).First(&curUser)
	if curUser.Username == user.Username {
		return errors.New("当前账号已存在")
	}

	hash, _ := encrypt.HashString(user.Password)
	user.Password = hash
	result := u.Create(user)
	if result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}

// Login 简单的登录
func (u *UserModel) Login(username string, password string) (*UserModel, error) {
	result := u.Where("username = ?", username).First(u)
	if result.RowsAffected == 0 {
		return nil, errors.New("用户不存在")
	}
	err := encrypt.CheckPassword(password, u.Password)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	return u, nil
}

// GetUserInfo GetInfo 获取用户数据
func (u *UserModel) GetUserInfo(id uint64) (*UserModel, error) {
	//sql := fmt.Sprintf("SELECT * FROM %s WHERE id = ? LIMIT 0,1", "user")
	//sql := "SELECT * FROM user WHERE id = ? LIMIT 0,1"
	mErr := u.Where("id = ?", id).First(u).Error

	if mErr != nil {
		if errors.Is(mErr, gorm.ErrRecordNotFound) {
			// 如果是记录未找到的错误，表示用户名不存在
			err := errors.New(fmt.Sprintf("未查询到相关数据[id:%s]", id))
			return nil, err
		}
	}
	return u, nil
}
