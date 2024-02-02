package authorization

import (
	"github.com/gin-gonic/gin"
	"sim/app/global/consts"
	"sim/app/services/user/token"
	"sim/app/util/response"
	"strings"
)

type HeaderParams struct {
	Authorization string `Header:"Authorization"`
}

// CheckLogin 登录校验中间件
func CheckLogin(c *gin.Context) {
	var header HeaderParams
	if err := c.ShouldBindHeader(&header); err != nil {
		response.TokenError(c, consts.TokenParseError)
		return
	}
	tokenArr := strings.Split(header.Authorization, " ")
	if len(tokenArr) != 2 {
		response.TokenError(c, consts.TokenInvalidError)
		return
	}
	user, code := token.CreateUserTokenFactory().IsEffective(tokenArr[1])
	if code != 1 {
		response.TokenError(c, code)
		return
	}
	c.Set("user", user)
	c.Next()
}
