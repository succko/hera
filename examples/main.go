package examples

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/succko/hera"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
)

func main() {
	defer hera.DeferHandle()
	// 注册模块
	modules := &hera.Modules{
		Db:        true,
		Redis:     true,
		Nacos:     true,
		Oss:       true,
		Flag:      true,
		Validator: true,
	}
	hera.RegisterModules(modules)

	// 注册任务
	hera.RegisterNacos(map[string]any{})
	hera.RegisterCron(func(c *cron.Cron) {
		// 定时任务列表
		_, _ = c.AddFunc("*/1 * * * *?", func() {
			zap.L().Info("cron task")
		})
	})
	hera.RegisterRocketMqConsumers(map[string]func(message []byte){})
	hera.RegisterMetaData([]func(){})
	hera.RegisterGrpc(func(server *grpc.Server) {
		// 注册 grpc 服务
	})
	// 注册 swagger
	hera.RegisterSwagger(func() {

	})

	// 启动服务器
	hera.RegisterRouter(func(router *gin.Engine) {
		// 注册 api 分组路由
		apiGroup := router.Group("/api")
		apiGroup.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
	})
	hera.RunHttpServer()
}
