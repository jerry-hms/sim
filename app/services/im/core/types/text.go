package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sim/app/services/im/core/interf"
)

type Text struct {
	Content string `json:"content"`
}

func (t *Text) ParseParams(msg map[string]interface{}) (interf.TypeInterface, error) {
	str, _ := json.Marshal(msg)
	err := json.Unmarshal(str, &t)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("message.%s 数据绑定失败", reflect.TypeOf(t).Name()))
	}
	return t, nil
}

func (t *Text) ParseContent() string {
	return t.Content
}
