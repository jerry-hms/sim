package types

import (
	"sim/app/services/im/core/interf"
)

type Image struct {
	Url    string `json:"url"`
	Height int64  `json:"height"`
	Width  int64  `json:"width"`
}

func (i *Image) ParseParams(msg map[string]interface{}) (interf.TypeInterface, error) {
	return DataBindToType(msg, i)
}
