## 初始化项目
```shell
go mod init hera-demo
go get -u github.com/succko/hera
```

## 编写main.go
```shell
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/succko/hera"
	"net/http"
)

func main() {
	defer hera.Hera.DeferHandle()
	modules := &hera.Modules{
		Nacos:    true,
		Db:       true,
		Redis:    true,
		Rocketmq: true,
	}
	_ = hera.Hera.Run(modules)

	r := gin.Default()

	// global.App 为注册的配置和服务

	// 测试路由
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 启动服务器
	r.Run(":8080")

}
```