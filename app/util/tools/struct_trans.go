package tools

import "encoding/json"

// InterfaceToStruct 空接口类型转换为对应的struct
func InterfaceToStruct(interf interface{}, stru interface{}) error {
	data, _ := json.Marshal(interf)
	err := json.Unmarshal(data, stru)
	if err != nil {
		return err
	}
	return nil
}
