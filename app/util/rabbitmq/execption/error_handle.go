package execption

import (
	"fmt"
)

// ErrorDeal 记录错误
func ErrorDeal(err error) error {
	if err != nil {
		//variable.ZapLog.Error(err.Error())
		fmt.Printf("rabbitmq报错: %s\n", err.Error())
	}
	return err
}
