package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/succko/hera/global"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetSwaggerGroupRoutes(router *gin.RouterGroup) {
	global.App.RunConfig.Swagger()
	router.GET("/*any", ginSwagger.DisablingWrapHandler(swaggerfiles.Handler, "SWAGGER"))
}
