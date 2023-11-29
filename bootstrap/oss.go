package bootstrap

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/succko/hera/global"
	"go.uber.org/zap"
)

func InitializeOss() *oss.Bucket {
	// 从环境变量中获取访问凭证。运行本代码示例之前，请确保已设置环境变量OSS_ACCESS_KEY_ID和OSS_ACCESS_KEY_SECRET。
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		zap.L().DPanic("初始化OSS失败", zap.Error(err))
	}

	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	client, err := oss.New(global.App.Config.Oss.Endpoint, global.App.Config.Oss.AccessKeyID, global.App.Config.Oss.AccessKeySecret, oss.SetCredentialsProvider(&provider))
	if err != nil {
		zap.L().DPanic("初始化OSS失败", zap.Error(err))
	}

	// 填写存储空间名称，例如examplebucket。
	bucket, err := client.Bucket(global.App.Config.Oss.BucketName)
	if err != nil {
		zap.L().DPanic("初始化OSS失败", zap.Error(err))
	}
	return bucket
}
