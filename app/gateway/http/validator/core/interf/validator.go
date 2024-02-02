package interf

import "github.com/gin-gonic/gin"

type ValidatorInterface interface {
	CheckParams(c *gin.Context)
}
