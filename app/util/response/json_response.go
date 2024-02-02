package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sim/app/global/consts"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ReturnJson(c *gin.Context, httpCode int, code int, msg string, data interface{}) {
	c.JSON(httpCode, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func Success(c *gin.Context, msg string, data interface{}) {
	ReturnJson(c, http.StatusOK, consts.RequestSuccess, msg, data)
}

func Fail(c *gin.Context, msg string, data interface{}) {
	ReturnJson(c, http.StatusInternalServerError, consts.RequestFail, msg, data)
	c.Abort()
}

func ValidatorFail(c *gin.Context, msg string) {
	ReturnJson(c, http.StatusInternalServerError, consts.RequestValidateFail, msg, nil)
	c.Abort()
}

func TokenError(c *gin.Context, code int) {
	var msg string
	switch code {
	case consts.TokenParseError:
		msg = consts.TokenParseErrorMsg
		break
	case consts.TokenInvalidError:
		msg = consts.TokenInvalidErrorMsg
		break
	case consts.TokenExpiredError:
		msg = consts.TokenExpiredErrorMsg
		break
	default:
		msg = consts.TokenParseErrorMsg
		break
	}
	ReturnJson(c, http.StatusInternalServerError, consts.TokenParseError, msg, nil)
	c.Abort()
}
