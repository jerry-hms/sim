package v1

import (
	"github.com/gin-gonic/gin"
	"sim/app/gateway/rpc"
	"sim/app/services/user/token"
	"sim/app/util/response"
	userPb "sim/idl/user"
)

type User struct {
}

// Register 注册
func (u *User) Register(c *gin.Context, userReq *userPb.UserRequest) {
	user, err := rpc.UserClient.UserRegister(c, userReq)
	if err != nil {
		response.Fail(c, err.Error(), nil)
		return
	}
	response.Success(c, "注册成功", user)
}

// Login 登录
func (u *User) Login(c *gin.Context) {
	userReq := &userPb.UserLoginRequest{
		Username: c.GetString("username"),
		Password: c.GetString("password"),
	}
	user, err := rpc.UserClient.UserLogin(c, userReq)
	if err != nil {
		response.Fail(c, err.Error(), nil)
		return
	}
	token, err := token.CreateUserTokenFactory().GenerateToken(user, 24*60*60)
	if err != nil {
		response.Fail(c, "token生成失败", nil)
		return
	}
	response.Success(c, "登录成功", gin.H{
		"token":     token,
		"user_info": user,
	})
}
