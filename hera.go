package hera

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/succko/hera/bootstrap"
	"github.com/succko/hera/global"
	"github.com/succko/hera/metadata"
	"github.com/xxl-job/xxl-job-executor-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type Modules struct {
	Db        bool
	Redis     bool
	Xxl       bool
	Nacos     bool
	Metadata  bool
	Rocketmq  bool
	Oss       bool
	Swagger   bool
	Grpc      bool
	Flag      bool
	Cron      bool
	Validator bool
}

var _modules = new(Modules)

// RegisterNacos 注册nacos配置
func RegisterNacos(m map[string]any) {
	_modules.Xxl = true
	global.App.RunConfig.Nacos = m
}

// RegisterCron 注册cron任务
func RegisterCron(f func(c *cron.Cron)) {
	_modules.Cron = true
	global.App.RunConfig.Cron = f
}

// RegisterRocketMqConsumers 注册rocketmq消费者
func RegisterRocketMqConsumers(m map[string]func(message []byte)) {
	_modules.Rocketmq = true
	global.App.RunConfig.RocketMqConsumers = m
}

// RegisterMetaData 注册元数据
func RegisterMetaData(fs []func(wg *sync.WaitGroup)) {
	_modules.Metadata = true
	global.App.RunConfig.MetaData = fs
}

// RegisterGrpc 注册grpc服务
func RegisterGrpc(f func(server *grpc.Server)) {
	_modules.Grpc = true
	global.App.RunConfig.Grpc = f
}

func RegisterRouter(f func(router *gin.Engine)) {
	global.App.RunConfig.Router = f
}

func RegisterXxl(f func(exec xxl.Executor)) {
	_modules.Xxl = true
	global.App.RunConfig.Xxl = f
}

func RegisterSwagger(f func()) {
	_modules.Swagger = true
	global.App.RunConfig.Swagger = f
}

func RegisterModules(modules *Modules) {
	_modules = modules
	//bootstrap.Modules = modules
}

// RunHttpServer 启动http服务
func RunHttpServer() {
	err := run()
	if err != nil {
		global.App.Log.Fatal("run http server error", zap.Error(err))
	}
	// 创建 TCP 监听器
	l, _ := net.Listen("tcp", ":"+global.App.Config.App.Port)
	bootstrap.RunHttpServer(l)
}

// RunGrpcServer 启动grpc服务
func RunGrpcServer() {
	err := run()
	if err != nil {
		global.App.Log.Fatal("run http server error", zap.Error(err))
	}
	bootstrap.RunGrpcServer()
}

func RunWsServer() {
	err := run()
	if err != nil {
		global.App.Log.Fatal("run http server error", zap.Error(err))
	}
	// 创建 TCP 监听器
	l, _ := net.Listen("tcp", ":"+global.App.Config.App.Port)
	bootstrap.RunWsServer(l)
}

// RunCMux 启动cmux服务
func RunCMux() {
	err := run()
	if err != nil {
		global.App.Log.Fatal("run http server error", zap.Error(err))
	}
	// 创建 TCP 监听器
	bootstrap.RunCMux()
}

func run() error {
	// 初始化配置
	if _, err := bootstrap.InitializeConfig(); err != nil {
		return err
	}

	// 初始化nacos配置
	if _modules.Nacos {
		if err := bootstrap.InitializeNacosConfig(); err != nil {
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
				metadata.Loader.InitializeMetadata()
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
				bootstrap.InitializeCron()
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
					global.App.RocketMqConsumers = bootstrap.InitializeRocketMqConsumers()
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

func DeferHandle() {
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
