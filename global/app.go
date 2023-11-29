package global

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"github.com/succko/hera/config"
	"github.com/xxl-job/xxl-job-executor-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
}

var App = new(app)
