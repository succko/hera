package config

type Rokcetmq struct {
	Addr      string     `mapstructure:"addr" json:"addr" yaml:"addr"`
	Consumers []Consumer ` yaml:"consumers"`
}

type Consumer struct {
	Topic string `mapstructure:"topic" json:"topic" yaml:"topic"`
	Func  func() `mapstructure:"func" json:"func" yaml:"func"`
}
