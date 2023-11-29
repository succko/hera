package config

type Nacos struct {
	Servers   []Server `mapstructure:"servers" json:"servers" yaml:"servers"`
	Namespace string   `mapstructure:"namespace" json:"namespace" yaml:"namespace"`
	Username  string   `mapstructure:"username" json:"username" yaml:"username"`
	Password  string   `mapstructure:"password" json:"password" yaml:"password"`
	DataId    string   `mapstructure:"data-id" json:"data-id" yaml:"data-id"`
}

type Server struct {
	ServerAddr string `mapstructure:"server-addr" json:"server-addr" yaml:"server-addr"`
	Port       uint64 `mapstructure:"port" json:"port" yaml:"port"`
}
