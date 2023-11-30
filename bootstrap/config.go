package bootstrap

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/succko/hera/global"
	"go.uber.org/zap"
	"os"
)

// InitializeConfig 初始化配置
func InitializeConfig() (*viper.Viper, error) {
	// 设置配置文件路径
	config := "config.yaml"
	// 生产环境可以通过设置环境变量来改变配置文件路径
	if configEnv := os.Getenv("VIPER_CONFIG"); configEnv != "" {
		config = configEnv
	}

	// 初始化 viper
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		zap.L().Error(fmt.Sprintf("read config failed: %s", err))
		return nil, err
	}

	// 监听配置文件
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.L().Info(fmt.Sprintf("config file changed: %s", in.Name))
		// 重载配置
		if err := v.Unmarshal(&global.App.Config); err != nil {
			zap.L().Error(fmt.Sprintf("config file error: %s", err))
		}
	})

	// 将配置赋值给全局变量
	if err := v.Unmarshal(&global.App.Config); err != nil {
		zap.L().Error(fmt.Sprintf("unmarshal config failed: %s", err))
		return nil, err
	}

	if global.App.Config.App.Env == gin.ReleaseMode {
		gin.SetMode(global.App.Config.App.Env)
	}

	return v, nil
}
