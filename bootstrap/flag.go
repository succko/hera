package bootstrap

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/succko/hera/global"
)

var (
	port = flag.String("port", "", "端口号: 默认以nacos配置为准")
	env  = flag.String("env", "debug", "环境: debug-测试 release-生产")
)

// InitializeFlag 初始化flag
func InitializeFlag() {
	flag.Parse()
	if *port != "" {
		global.App.Config.App.Port = *port
	}
	if *env == gin.ReleaseMode {
		global.App.Config.App.Env = *env
		gin.SetMode(*env)
	}
}
