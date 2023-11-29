package config

type Oss struct {
	AccessKeyID     string `mapstructure:"access_key_id" json:"access_key_id" yaml:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret" json:"access_key_secret" yaml:"access_key_secret"`
	OssSessionToken string `mapstructure:"oss_session_token" json:"oss_session_token" yaml:"oss_session_token"`
	Endpoint        string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	BucketName      string `mapstructure:"bucket_name" json:"bucket_name" yaml:"bucket_name"`
}
