package global

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"github.com/succko/hera/config"
	"github.com/xxl-job/xxl-job-executor-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"sync"
)

type app struct {
	ConfigViper       *viper.Viper
	Config            config.Configuration
	Log               *zap.Logger
	DB                *gorm.DB
	Redis             *redis.Client
	Xxl               xxl.Executor
	Oss               *oss.Bucket
	RocketMqProducer  rocketmq.Producer
	RocketMqConsumers []rocketmq.PushConsumer
	RunConfig         RunConfig
}

type RunConfig struct {
	Nacos             map[string]any
	Cron              func(c *cron.Cron)
	RocketMqConsumers map[string]func(message []byte)
	MetaData          []func(wg *sync.WaitGroup)
	Grpc              func(server *grpc.Server)
	Xxl               func(exec xxl.Executor)
	Router            func(router *gin.Engine)
	Swagger           func()
}

var App = new(app)
