package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/succko/hera"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

func main() {
	defer hera.DeferHandle()
	// 注册模块
	modules := &hera.Modules{
		Db:        true,
		Redis:     true,
		Xxl:       true,
		Nacos:     true,
		Metadata:  true,
		Rocketmq:  true,
		Oss:       true,
		Grpc:      true,
		Flag:      true,
		Cron:      true,
		Validator: true,
	}
	hera.RegisterModules(modules)

	// 注册任务
	tasks := &hera.Tasks{
		Nacos: map[string]any{},
		Cron: func(c *cron.Cron) {
			// 定时任务列表
			_, _ = c.AddFunc("*/1 * * * *?", func() {
				zap.L().Info("cron task")
			})
		},
		RocketMqConsumers: map[string]func(message []byte){},
		MetaData:          []func(wg *sync.WaitGroup){},
	}
	hera.RegisterTasks(tasks)

	// 启动服务器
	hera.RunServer(func(router *gin.Engine) {
		// 注册 api 分组路由
		apiGroup := router.Group("/api")
		apiGroup.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
	})
}
