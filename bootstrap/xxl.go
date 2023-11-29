package bootstrap

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/succko/hera/global"
	"github.com/xxl-job/xxl-job-executor-go"
	"log"
)

func InitializeXxl() xxl.Executor {
	exec := xxl.NewExecutor(
		xxl.ServerAddr(global.App.Config.Xxl.ServerAddr),
		xxl.AccessToken(global.App.Config.Xxl.AccessToken), //请求令牌(默认为空)
		//xxl.ExecutorIp(global.App.Config.Xxl.ExecutorIp),     //可自动获取
		xxl.ExecutorPort(global.App.Config.App.Port),   //默认9999（非必填）
		xxl.RegistryKey(global.App.Config.App.AppName), //执行器名称
		xxl.SetLogger(&xxlLogger{}),                    //自定义日志
	)
	exec.Init()
	//设置日志查看handler
	exec.LogHandler(func(req *xxl.LogReq) *xxl.LogRes {
		return &xxl.LogRes{Code: 200, Msg: "", Content: xxl.LogResContent{
			FromLineNum: req.FromLineNum,
			ToLineNum:   2,
			LogContent:  "这个是自定义日志handler",
			IsEnd:       true,
		}}
	})
	//defer exec.Stop()
	//log.Fatal(exec.Run())
	return exec
}

func XxlJobMux(e *gin.Engine, exec xxl.Executor) {
	//注册的gin的路由
	e.POST("run", gin.WrapF(exec.RunTask))
	e.POST("kill", gin.WrapF(exec.KillTask))
	e.POST("log", gin.WrapF(exec.TaskLog))
	e.POST("beat", gin.WrapF(exec.Beat))
	e.POST("idleBeat", gin.WrapF(exec.IdleBeat))
	//注册任务handler
	if global.App.RunConfig.Xxl != nil {
		global.App.RunConfig.Xxl(exec)
	}
}

// xxl.Logger接口实现
type xxlLogger struct{}

func (l *xxlLogger) Info(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf("xxl-job info - "+format, a...))
}

func (l *xxlLogger) Error(format string, a ...interface{}) {
	log.Println(fmt.Sprintf("xxl-job error - "+format, a...))
}
