package bootstrap

import (
	"log"
	"os"
	api_register_validator "sim/app/gateway/http/validator/common/register_validator"
	my_errors "sim/app/global/error"
	"sim/app/global/variable"
	"sim/app/providers/queue"
	"sim/app/services/sys_log"
	"sim/app/services/websocket/client"
	"sim/app/util/websocket"
	"sim/app/util/yml_config"
	"sim/app/util/zap"
)

func checkRequiredFolders() {
	//1.检查配置文件是否存在
	if _, err := os.Stat(variable.BasePath + "/config/config.yml"); err != nil {
		log.Fatal(my_errors.ErrorsConfigYamlNotExists + err.Error())
	}
	//2.检查public目录是否存在
	if _, err := os.Stat(variable.BasePath + "/public/"); err != nil {
		log.Fatal(my_errors.ErrorsPublicNotExists + err.Error())
	}
	//3.检查storage/logs 目录是否存在
	if _, err := os.Stat(variable.BasePath + "/storage/logs/"); err != nil {
		log.Fatal(my_errors.ErrorsStorageLogsNotExists + err.Error())
	}
	// 4.自动创建软连接、更好的管理静态资源
	if _, err := os.Stat(variable.BasePath + "/public/storage"); err == nil {
		if err = os.RemoveAll(variable.BasePath + "/public/storage"); err != nil {
			log.Fatal(my_errors.ErrorsSoftLinkDeleteFail + err.Error())
		}
	}
	if err := os.Symlink(variable.BasePath+"/storage/app", variable.BasePath+"/public/storage"); err != nil {
		log.Fatal(my_errors.ErrorsSoftLinkCreateFail + err.Error())
	}
}

func init() {
	checkRequiredFolders()
	// 注册api的validator验证器
	api_register_validator.Handler()

	variable.ConfigYml = yml_config.CreateYamlFactory()
	variable.ConfigYml.ConfigFileChangeListen()

	// 初始化全局日志句柄，并载入日志钩子处理函数
	variable.ZapLog = zap.CreateZapFactory(sys_log.ZapLogHandler)

	// 开启websocket服务
	if variable.ConfigYml.GetInt("websocket.Start") == 1 {

		variable.WebsocketHub = websocket.CreateHubFactory()
		// 启动客户端管理器
		variable.WebsocketManage = client.CreateManage()
		if Wh, ok := variable.WebsocketHub.(*websocket.Hub); ok {
			go Wh.Run()
		}
	}
	// 将队列处理handle注入容器
	queue.RegisterQueue()
}
