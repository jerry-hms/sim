package data_transfer

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"sim/app/gateway/http/validator/core/interf"
)

func DataAddContext(validatorInterface interf.ValidatorInterface, dataPrefix string, c *gin.Context) *gin.Context {
	var tempJson interface{}
	if data, err := json.Marshal(validatorInterface); err == nil {
		if err2 := json.Unmarshal(data, &tempJson); err2 == nil {
			if value, ok := tempJson.(map[string]interface{}); ok {
				for k, v := range value {
					c.Set(dataPrefix+k, v)
				}
				return c
			}
		}
	}

	return nil
}
