package hera

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/succko/hera/bootstrap"
	"github.com/succko/hera/global"
	"github.com/succko/hera/metadata"
	"go.uber.org/zap"
	"sync"
)

type h struct {
}

var Hera = new(h)

func main() {
	defer Hera.DeferHandle()
	mod := &Modules{
		Db:        true,
		Redis:     true,
		Xxl:       true,
		Nacos:     true,
		Metadata:  true,
		Rocketmq:  true,
		Oss:       true,
		Grpc:      true,
		Flag:      true,
		Cron:      true,
		Validator: true,
	}
	if err := Hera.Run(mod); err != nil {
		zap.L().DPanic("Initialization failed", zap.Error(err))
	}
	fmt.Println("启动成功")
}

type Tasks struct {
	Nacos             map[string]any
	Cron              func(c *cron.Cron)
	RocketMqConsumers map[string]func(message []byte)
	MetaData          []func(wg *sync.WaitGroup)
}

var _tasks *Tasks

func (h *h) Register(tasks *Tasks) *h {
	_tasks = tasks
	return h
}

type Modules struct {
	Db        bool
	Redis     bool
	Xxl       bool
	Nacos     bool
	Metadata  bool
	Rocketmq  bool
	Oss       bool
	Grpc      bool
	Flag      bool
	Cron      bool
	Validator bool
}

var _modules *Modules

func (h *h) Run(modules *Modules) error {
	_modules = modules
	if _tasks == nil {
		h.Register(&Tasks{
			Nacos:             map[string]any{},
			Cron:              func(c *cron.Cron) {},
			RocketMqConsumers: map[string]func(message []byte){},
			MetaData:          []func(wg *sync.WaitGroup){},
		})
	}

	// 初始化配置
	if _, err := bootstrap.InitializeConfig(); err != nil {
		return err
	}

	// 初始化nacos配置
	if _modules.Nacos {
		if err := bootstrap.InitializeNacosConfig(_tasks.Nacos); err != nil {
			return err
		}
	}

	// 初始化flag
	if _modules.Flag {
		bootstrap.InitializeFlag()
	}

	// 初始化日志
	global.App.Log = bootstrap.InitializeLog()

	// 初始化数据库
	if _modules.Db {
		global.App.DB = bootstrap.InitializeDB()
	}

	var wg sync.WaitGroup

	inits := make([]func() error, 0)
	if _modules.Db {
		inits = append(inits, // 初始化验证器
			func() error {
				defer wg.Done()
				return bootstrap.InitializeValidator()
			})
	}

	if _modules.Redis {
		inits = append(inits, // 初始化Redis
			func() error {
				defer wg.Done()
				global.App.Redis = bootstrap.InitializeRedis()
				return nil
			})
	}

	if _modules.Xxl {
		inits = append(inits, // 初始化Xxl
			func() error {
				defer wg.Done()
				global.App.Xxl = bootstrap.InitializeXxl()
				return nil
			})
	}

	if _modules.Metadata {
		inits = append(inits, // 初始化元数据
			func() error {
				defer wg.Done()
				metadata.Loader.InitializeMetadata(_tasks.MetaData)
				return nil
			})
	}

	if _modules.Oss {
		inits = append(inits, // 初始化OSS
			func() error {
				defer wg.Done()
				global.App.Oss = bootstrap.InitializeOss()
				return nil
			})
	}

	if _modules.Cron {
		inits = append(inits, // 初始化Cron
			func() error {
				defer wg.Done()
				bootstrap.InitializeCron(_tasks.Cron)
				return nil
			})
	}

	if _modules.Rocketmq {
		inits = append(inits, // 初始化RocketMq
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
					global.App.RocketMqConsumers = bootstrap.InitializeRocketMqConsumers(_tasks.RocketMqConsumers)
				}()
				w.Wait()
				return nil
			})
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

func (h *h) DeferHandle() {
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