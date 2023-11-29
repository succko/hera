package config

type Xxl struct {
	ServerAddr  string `mapstructure:"server_addr" json:"server_addr" yaml:"server_addr"`
	AccessToken string `mapstructure:"access_token" json:"access_token" yaml:"access_token"`
	ExecutorIp  string `mapstructure:"executor_ip" json:"executor_ip" yaml:"executor_ip"`
	//ExecutorPort string `mapstructure:"executor_port" json:"executor_port" yaml:"executor_port"`
}
