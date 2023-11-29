package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/succko/hera/config"
)

type app struct {
	//ConfigViper       *viper.Viper
	Config config.Configuration
	//Log               *zap.Logger
	//DB                *gorm.DB
	Redis *redis.Client
	//Xxl               xxl.Executor
	//Oss               *oss.Bucket
	//RocketMqProducer  rocketmq.Producer
	//RocketMqConsumers []rocketmq.PushConsumer
}

var App = new(app)
