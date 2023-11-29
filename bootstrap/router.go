package bootstrap

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/succko/hera/global"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func setupRouter(f func(router *gin.Engine)) *gin.Engine {
	router := gin.Default()
	f(router)
	return router
}

// RunServer 启动服务器
func RunServer(f func(router *gin.Engine)) {
	r := setupRouter(f)

	srv := &http.Server{
		Addr:    ":" + global.App.Config.App.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error("listen: ", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown:", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
