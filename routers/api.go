package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// SetApiGroupRouters 定义 api 分组路由
func SetApiGroupRouters(router *gin.RouterGroup) {
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}
