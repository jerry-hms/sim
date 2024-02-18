package api

import (
	"github.com/gin-gonic/gin"
	"sim/app/util/jwt"
	userPb "sim/idl/pb/user"
)

// GetLoginUser 获取登录的用户信息
func GetLoginUser(c *gin.Context) userPb.UserResponse {
	user, _ := c.Get("user")
	claims := user.(*jwt.Claims)
	return claims.Info
}
