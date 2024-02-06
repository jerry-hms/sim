package response

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"sim/app/global/consts"
	my_errors "sim/app/global/error"
	"sim/app/util/validator_translation"
	"strings"
)

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

func RouterNotFoundFail(c *gin.Context) {
	ReturnJson(c, http.StatusBadRequest, consts.RouterNotFoundFail, consts.RouterNotFoundFailMsg, nil)
	c.Abort()
}

// ErrorSystem 系统执行代码错误
func ErrorSystem(c *gin.Context, msg string, data interface{}) {
	ReturnJson(c, http.StatusInternalServerError, consts.SystemError, consts.SystemErrorMsg+msg, data)
	c.Abort()
}

// ValidatorError 翻译表单参数验证器出现的校验错误
func ValidatorError(c *gin.Context, validtorIns interface{}, err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		wrongParam := validator_translation.RemoveTopStruct(validtorIns, errs)
		ReturnJson(c, http.StatusBadRequest, consts.ValidatorCheckParamsFail, consts.ValidatorCheckParamsFailMsg, wrongParam)
	} else {
		errStr := err.Error()
		// multipart:nextpart:eof 错误表示验证器需要一些参数，但是调用者没有提交任何参数
		if strings.ReplaceAll(strings.ToLower(errStr), " ", "") == "multipart:nextpart:eof" {
			ReturnJson(c, http.StatusBadRequest, consts.ValidatorCheckParamsFail, consts.ValidatorCheckParamsFailMsg, gin.H{"tips": my_errors.ErrorNotAllParamsIsBlank})
		} else {
			ReturnJson(c, http.StatusBadRequest, consts.ValidatorCheckParamsFail, consts.ValidatorCheckParamsFailMsg, gin.H{"tips": errStr})
		}
	}
	c.Abort()
}

// TokenError 抛出token错误
func TokenError(c *gin.Context, code int) {
	var msg string
	switch code {
	case consts.TokenParseError:
		msg = consts.TokenParseErrorMsg
	case consts.TokenInvalidError:
		msg = consts.TokenInvalidErrorMsg
	case consts.TokenExpiredError:
		msg = consts.TokenExpiredErrorMsg
	case consts.TokenFormatError:
		msg = consts.TokenFormatErrorMsg
	default:
		msg = consts.TokenParseErrorMsg
	}
	ReturnJson(c, http.StatusInternalServerError, consts.TokenParseError, msg, nil)
	c.Abort()
}
