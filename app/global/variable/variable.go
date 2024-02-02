package variable

import (
	"go.uber.org/zap"
	"log"
	"os"
	my_errors "sim/app/global/error"
	"sim/app/util/yml_config/interf"
	"strings"
)

// ginskeleton 封装的全局变量全部支持并发安全，请放心使用即可
// 开发者自行封装的全局变量，请做好并发安全检查与确认

var (
	BasePath           string                  // 定义项目的根目录
	EventDestroyPrefix = "Destroy_"            //  程序退出时需要销毁的事件前缀
	ConfigKeyPrefix    = "Config_"             //  配置文件键值缓存时，键的前缀
	DateFormat         = "2006-01-02 15:04:05" //  设置全局日期时间格式

	// 全局日志指针
	ZapLog *zap.Logger
	// 全局配置文件
	ConfigYml interf.YmlConfigInterf // 全局配置文件指针
	//ConfigGormv2Yml interf.YmlConfigInterf // 全局配置文件指针

	//websocket
	WebsocketHub              interface{}
	WebsocketHandshakeSuccess = `{"code":200,"msg":"ws连接成功","data": {"client_id": "%s"}}`
	WebsocketServerPingMsg    = "Server->Ping->Client"

	//  用户自行定义其他全局变量 ↓

)

func init() {
	// 1.初始化程序根目录
	if curPath, err := os.Getwd(); err == nil {
		// 路径进行处理，兼容单元测试程序程序启动时的奇怪路径
		if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-study_gorm") {
			BasePath = strings.Replace(strings.Replace(curPath, `\test`, "", 1), `/study_gorm`, "", 1)
		} else {
			BasePath = curPath
		}
	} else {
		log.Fatal(my_errors.ErrorsBasePath)
	}
}
