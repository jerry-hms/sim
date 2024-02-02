package bootstrap

import (
	"log"
	"os"
	api_register_validator "sim/app/gateway/http/validator/common/register_validator"
	my_errors "sim/app/global/error"
	"sim/app/global/variable"
	"sim/app/util/websocket"
	"sim/app/util/yml_config"
)

func checkRequiredFolders() {
	if _, err := os.Stat(variable.BasePath + "/config/config.yml"); err != nil {
		log.Fatal(my_errors.ErrorsConfigYamlNotExists + err.Error())
	}
}

func init() {
	checkRequiredFolders()
	// 注册api的validator验证器
	api_register_validator.Handler()

	variable.ConfigYml = yml_config.CreateYamlFactory()
	variable.ConfigYml.ConfigFileChangeListen()

	// 开启websocket服务
	if variable.ConfigYml.GetInt("websocket.Start") == 1 {
		variable.WebsocketHub = websocket.CreateHubFactory()
		if Wh, ok := variable.WebsocketHub.(*websocket.Hub); ok {
			go Wh.Run()
		}
	}
}
