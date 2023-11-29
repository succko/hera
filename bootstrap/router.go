package bootstrap

import (
	"context"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"github.com/succko/hera/global"
	"github.com/succko/hera/routes"
	"github.com/succko/hera/ws"
	"go.uber.org/zap"
	"golang.org/x/net/websocket"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// 设置路由
func setupRouter() *gin.Engine {
	if global.App.Config.App.Env == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	// 静态文件 缓存测试
	inMemoryStore := persistence.NewInMemoryStore(60 * time.Second)
	r.GET("/cache_ping", cache.CachePage(inMemoryStore, time.Minute, func(c *gin.Context) {
		r.StaticFile("/test", "./static/test.json")
	}))

	// 使用自定义的日志和恢复中间件
	//r.Use(gin.Logger(), gin.Recovery())
	r.Use(GinLogger(), GinRecovery(true))

	// 注册 ping 路由
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	// 注册 ws 路由
	r.GET("/ws", func(ctx *gin.Context) {
		ws.ServeWs(ctx.Writer, ctx.Request)
	})

	// 注册 api 分组路由
	global.App.RunConfig.Router(r)

	// 注册 swagger 分组路由

	if gin.Mode() != gin.ReleaseMode {
		swaggerGroup := r.Group("/swagger")
		routes.SetSwaggerGroupRoutes(swaggerGroup)
	}
	return r
}

// RunCMux 运行 CMux
func RunCMux() {
	// 创建 TCP 监听器
	l, err := net.Listen("tcp", ":"+global.App.Config.App.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 创建 CMux 实例
	m := cmux.New(l)
	// 创建 gRPC 匹配规则
	//grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	// 创建 ws 匹配规则
	//wsL := m.Match(cmux.HTTP1HeaderField("Upgrade", "websocket"))
	// 否则使用 HTTP 匹配规则
	httpL := m.Match(cmux.Any())

	// 启动 gRPC 服务器
	//go RunGrpcServer(grpcL)
	go RunGrpcServer()
	//go RunWsServer(wsL)
	// 启动 WebSocket 服务器
	go ws.SingletonHub().Run()
	// 启动 HTTP 服务器
	go RunHttpServer(httpL)

	// 启动 CMux
	if err := m.Serve(); !strings.Contains(err.Error(), "use of closed network connection") {
		zap.L().DPanic("CMux Serve error", zap.Error(err))
	}
}

// RunGrpcServer 运行 gRPC 服务器
func RunGrpcServer() {
	// 创建 gRPC 服务器实例
	server := grpc.NewServer()
	// 注册服务
	global.App.RunConfig.Grpc(server)
	//pb.RegisterUserServer(s, UserServer)
	//pb.RegisterTestServer(s, TestServer)
	//pb.RegisterPalaServer(s, PalaServer)

	// 启动 gRPC 服务器
	port, _ := strconv.Atoi(global.App.Config.App.Port)
	lis, _ := net.Listen("tcp", ":"+strconv.Itoa(port+10000))
	if err := server.Serve(lis); err != cmux.ErrListenerClosed {
		zap.L().DPanic("gRPC server error", zap.Error(err))
	}
}

// RunWsServer 运行 WebSocket 服务器
func RunWsServer(l net.Listener) {
	// 创建 WebSocket 服务器
	server := &http.Server{Handler: websocket.Handler(func(conn *websocket.Conn) {
		if _, err := io.Copy(conn, conn); err != nil {
			panic(err)
		}
	})}
	// 启动 WebSocket 服务器
	if err := server.Serve(l); err != cmux.ErrListenerClosed {
		zap.L().Error("WebSocket server error", zap.Error(err))
	}
}

// RunHttpServer 运行 HTTP 服务器
func RunHttpServer(l net.Listener) {
	// 设置路由
	r := setupRouter()

	// 添加 XxlJob 路由
	//if Modules.Xxl {
	XxlJobMux(r, global.App.Xxl)
	//}

	// 添加 pprof 性能分析 路由
	pprof.Register(r)

	// 创建 HTTP 服务器实例
	s := &http.Server{
		Handler: r,
	}

	// 启动 HTTP 服务器
	go func() {
		if err := s.Serve(l); err != nil && err != cmux.ErrListenerClosed {
			zap.L().DPanic("HTTP server error", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown:", zap.Error(err))
	}
	select {
	case <-ctx.Done():
		zap.L().Info("timeout of 5 seconds.")
	}
	zap.L().Info("Server exiting")
}
