package user

import (
	"context"
	"sim/app/model"
	pb "sim/idl/pb/user"
	"sync"
)

var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

type UserSrv struct {
	pb.UnimplementedUserServiceServer
	NodeAddr string
}

// GetUserSrv 使用sync.Once实现单例
func GetUserSrv(addr string) *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
		UserSrvIns.NodeAddr = addr
	})
	return UserSrvIns
}

// UserRegister 注册用户
func (u *UserSrv) UserRegister(ctx context.Context, req *pb.UserRequest) (resp *pb.UserResponse, err error) {
	user := &model.User{
		Username: req.UserName,
		Password: req.Password,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Mobile:   req.Mobile,
	}
	err = model.CreateUserFactory().Register(user)
	if err != nil {
		return
	}

	resp = toResponse(user)
	return
}

// UserLogin 用户登录
func (u *UserSrv) UserLogin(ctx context.Context, req *pb.UserLoginRequest) (resp *pb.UserResponse, err error) {
	user, err := model.CreateUserFactory().Login(req.Username, req.Password)
	if err != nil {
		return
	}

	resp = toResponse(user)
	return
}

func (u *UserSrv) UserInfo(ctx context.Context, req *pb.UserInfoRequest) (*pb.UserResponse, error) {
	info, err := model.CreateUserFactory().GetUserInfo(req.Id)
	if err != nil {
		return nil, err
	}

	return toResponse(info), nil
}

// 转换为响应数据
func toResponse(user *model.User) *pb.UserResponse {
	resp := &pb.UserResponse{
		Id:        user.Id,
		NickName:  user.Nickname,
		Avatar:    user.Avatar,
		Mobile:    user.Mobile,
		CreatedAt: user.BaseModel.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp
}
