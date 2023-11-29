package config

type Configuration struct {
	//App            App      `mapstructure:"app" json:"app" yaml:"app"`
	//Log            Log      `mapstructure:"log" json:"log" yaml:"log"`
	//Database       Database `mapstructure:"database" json:"database" yaml:"database"`
	Redis Redis `mapstructure:"redis" json:"redis" yaml:"redis"`
	//Jwt            Jwt      `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	//Xxl            Xxl      `mapstructure:"xxl" json:"xxl" yaml:"xxl"`
	//Nacos          Nacos    `mapstructure:"nacos" json:"nacos" yaml:"nacos"`
	//Rokcetmq       Rokcetmq `mapstructure:"rokcetmq" json:"rokcetmq" yaml:"rokcetmq"`
	//Oss            Oss      `mapstructure:"oss" json:"oss" yaml:"oss"`
	//UpdateVersion  UpdateVersion
	//StartUpIos     StartUpIos
	//StartUpAndroid StartUpAndroid
}
