package hera

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/succko/hera/bootstrap"
	"github.com/succko/hera/global"
	"go.uber.org/zap"
	"sync"
)

// @title Golang项目
// @version 1.0.0
// @description 这是一个Golang编写的经典项目
// @termsOfService http://swagger.io/terms/

// @contact.name pala
// @contact.url http://www.swagger.io/support

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:48085
func main() {
	defer deferHandle()
	if err := Init(); err != nil {
		zap.L().DPanic("Initialization failed", zap.Error(err))
	}
	fmt.Println("启动成功")
}

func Init() error {
	// 初始化配置
	if _, err := bootstrap.InitializeConfig(); err != nil {
		return err
	}

	// 初始化nacos配置
	if err := bootstrap.InitializeNacosConfig(); err != nil {
		return err
	}

	// 初始化flag
	bootstrap.InitializeFlag()

	// 初始化日志
	global.App.Log = bootstrap.InitializeLog()

	// 初始化数据库
	global.App.DB = bootstrap.InitializeDB()

	var wg sync.WaitGroup

	inits := []func() error{
		// 初始化验证器
		func() error {
			defer wg.Done()
			return bootstrap.InitializeValidator()
		},
		// 初始化Redis
		func() error {
			defer wg.Done()
			bootstrap.InitializeRedis()
			return nil
		},
		// 初始化Xxl
		func() error {
			defer wg.Done()
			global.App.Xxl = bootstrap.InitializeXxl()
			return nil
		},
		// 初始化元数据
		func() error {
			defer wg.Done()
			//metadata.Loader.InitializeMetadata()
			return nil
		},
		// 初始化OSS
		func() error {
			defer wg.Done()
			global.App.Oss = bootstrap.InitializeOss()
			return nil
		},
		func() error {
			defer wg.Done()
			bootstrap.InitializeCron(func(c *cron.Cron) {
				// TODO
			})
			return nil
		},
		// 初始化RocketMq
		func() error {
			defer wg.Done()
			var w sync.WaitGroup
			w.Add(2)
			go func() {
				defer w.Done()
				global.App.RocketMqProducer = bootstrap.InitializeRocketMqProducer()
			}()
			go func() {
				defer w.Done()
				// TODO
				global.App.RocketMqConsumers = bootstrap.InitializeRocketMqConsumers(map[string]func(message []byte){})
			}()
			w.Wait()
			return nil
		},
	}

	wg.Add(len(inits))

	for _, f := range inits {
		go func(initFunc func() error) {
			if err := initFunc(); err != nil {
				zap.L().Error("Initialization task failed", zap.Error(err))
			}
		}(f)
	}

	// 等待所有初始化任务完成
	wg.Wait()

	return nil
}

func deferHandle() {
	zap.L().Info("defer handle trigger")

	// 程序关闭前，释放数据库连接
	if global.App.DB != nil {
		db, _ := global.App.DB.DB()
		if err := db.Close(); err == nil {
			zap.L().Info("defer db close success")
		} else {
			zap.L().Error("defer db close error", zap.Error(err))
		}
	}

	if global.App.RocketMqProducer != nil {
		if err := global.App.RocketMqProducer.Shutdown(); err == nil {
			zap.L().Info("defer mq producer shutdown success")
		} else {
			zap.L().Error("defer mq producer shutdown error", zap.Error(err))
		}
	}

	if global.App.RocketMqConsumers != nil && len(global.App.RocketMqConsumers) > 0 {
		for _, consumer := range global.App.RocketMqConsumers {
			if err := consumer.Shutdown(); err == nil {
				zap.L().Info("defer mq consumer shutdown success")
			} else {
				zap.L().Error("defer mq consumer shutdown error", zap.Error(err))
			}
		}
	}

	zap.L().Info("defer handle end")
}
