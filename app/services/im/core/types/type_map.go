package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sim/app/services/im/core/interf"
)

const (
	TextMessage  string = "text"
	ImageMessage string = "image"
)

type messageTypeMap map[string]interf.TypeInterface

var MessageTypes = messageTypeMap{
	TextMessage:  &Text{},
	ImageMessage: &Image{},
}

// DataBindToType 数据绑定到类型
func DataBindToType(data map[string]interface{}, x interf.TypeInterface) (interf.TypeInterface, error) {
	str, _ := json.Marshal(data)
	err := json.Unmarshal(str, &x)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("message.%s 数据绑定失败", reflect.TypeOf(x).Name()))
	}
	return x, nil
}
