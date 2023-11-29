package bootstrap

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/succko/hera/global"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func InitializeNacosConfig(m map[string]any) error {
	// 创建clientConfig的另一种方式
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(global.App.Config.Nacos.Namespace), //当namespace是public时，此处填空字符串。
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
		constant.WithUsername(global.App.Config.Nacos.Username),
		constant.WithPassword(global.App.Config.Nacos.Password),
	)

	// 创建serverConfig的另一种方式
	serverConfigs := make([]constant.ServerConfig, len(global.App.Config.Nacos.Servers))
	for i, server := range global.App.Config.Nacos.Servers {
		serverConfigs[i] = *constant.NewServerConfig(
			server.ServerAddr,
			server.Port,
			constant.WithScheme("http"),
			constant.WithContextPath("/nacos"),
		)
	}

	m[global.App.Config.App.AppName+"-"+gin.Mode()+".yaml"] = &global.App.Config
	wg.Add(len(m))
	for k, v := range m {
		k := k
		v := v
		go func() {
			defer wg.Done()
			err := listenConfig(k, v, clientConfig, serverConfigs)
			if err != nil {
				zap.L().Error("listen config error", zap.Error(err))
			}
			zap.L().Info("listen config", zap.String("dataId", k))
		}()
	}
	wg.Wait()
	zap.L().Info("nacos config initialized")
	return nil
}

func listenConfig[T any](DataId string, t T, clientConfig constant.ClientConfig, serverConfigs []constant.ServerConfig) error {
	// 创建动态配置客户端的另一种方式 (推荐)
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		zap.L().Panic("初始化配置失败", zap.Error(err))
	}
	Group := "DEFAULT_GROUP"
	// 将配置赋值给全局变量
	param := vo.ConfigParam{DataId: DataId, Group: Group}
	content, err := configClient.GetConfig(param)
	if strings.HasSuffix(DataId, ".yaml") {
		parseYaml(content, t)
	} else {
		parseJson(content, t)
	}
	// 监听配置文件
	param.OnChange = func(namespace, group, dataId, data string) {
		// 重载配置
		if strings.HasSuffix(DataId, ".yaml") {
			// 脱敏日志
			zap.L().Info("config file changed, group:" + group + ", dataId:" + dataId)
			parseYaml(content, t)
		} else {
			zap.L().Info("config file changed, group:" + group + ", dataId:" + dataId + ", data:" + data)
			parseJson(content, t)
		}
	}
	err = configClient.ListenConfig(param)
	return err
}

func parseYaml[T any](data string, t T) {
	if err := yaml.Unmarshal([]byte(data), t); err != nil {
		zap.L().Error("初始化配置失败", zap.Error(err))
		return
	}
	zap.L().Info("初始化配置成功", zap.Any("config", global.App.Config))
	// 处理OSS环境变量
	_ = os.Setenv("OSS_ACCESS_KEY_ID", global.App.Config.Oss.AccessKeyID)
	_ = os.Setenv("OSS_ACCESS_KEY_SECRET", global.App.Config.Oss.AccessKeySecret)
	_ = os.Setenv("OSS_SESSION_TOKEN", global.App.Config.Oss.OssSessionToken)
}

func parseJson[T any](data string, t T) {
	if err := json.Unmarshal([]byte(data), t); err != nil {
		zap.L().Error("初始化配置失败", zap.Error(err))
		return
	}
	zap.L().Info("初始化配置成功", zap.Any("config", t))
}
