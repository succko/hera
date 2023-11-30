package config

type Modules struct {
	Db        bool
	Redis     bool
	Nacos     bool
	Oss       bool
	Flag      bool
	Validator bool
	Ws        bool
	Http      bool
}

type AllModules struct {
	Modules
	Xxl      bool
	Metadata bool
	Rocketmq bool
	Swagger  bool
	Grpc     bool
	Cron     bool
}
