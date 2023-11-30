## 初始化项目
```shell
go mod init hera-demo
go get -u github.com/succko/hera
```

## 编写main.go
```go
package main

import (
	"github.com/succko/hera"
	"github.com/succko/hera/config"
)

func main() {
	defer hera.DeferHandle()
	// 注册模块
	hera.RegisterModules(&config.Modules{
		Db:        true,
		Redis:     true,
		Nacos:     true,
		Oss:       true,
		Flag:      true,
		Validator: true,
	})
	// cron
	hera.RegisterCron(task.RegisterCron)
	// xxl
	hera.RegisterXxl(task.RegisterXxl)
	// 消息队列消费者
	hera.RegisterRocketMqConsumers(mq.RegisterRocketMqConsumers)
	// 元数据加载
	hera.RegisterMetaData(metadata.RegisterMetaData)
	// grpc
	hera.RegisterGrpc(server.RegisterGrpc)
	// 路由
	hera.RegisterRouter(routes.RegisterRouter)
	// swagger
	hera.RegisterSwagger(routes.RegisterSwagger)
	// 启动服务
	hera.RunCMux(true, true, true)
}
```

## config.yaml
```yaml
app:
  app_name: hera-demo       # 应用名称
nacos:
  servers:
    - server-addr: 
      port: 
  namespace: 
  username: 
  password: 

```