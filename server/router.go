package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"proxy-engineering-thesis/server/controller"
)

func NewRouter(proxyController *controller.ProxyController, sourceController *controller.DataSourceController) *gin.Engine {
	service := gin.Default()

	service.GET("", func(context *gin.Context) {
		context.JSON(http.StatusOK, "welcome home")
	})

	service.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	service.Group("/").StaticFile("/favicon.ico", "./resources/favicon.ico")

	router := service.Group("/api")

	proxyRouter := router.Group("/proxy")
	proxyRouter.GET("", proxyController.GetAll)
	proxyRouter.GET("/:id", proxyController.FindById)
	proxyRouter.GET("/:id/sessions/count", proxyController.GetProxySessionsCount)
	proxyRouter.POST("", proxyController.Create)
	proxyRouter.PUT("/:id/start", proxyController.StartProxy)
	proxyRouter.PUT("/:id/stop", proxyController.StopProxy)
	proxyRouter.DELETE("/:id", proxyController.Delete)

	dataSourceRouter := router.Group("/datasource")
	dataSourceRouter.GET("", sourceController.GetAll)
	dataSourceRouter.POST("", sourceController.Create)
	dataSourceRouter.GET("/:id", sourceController.FindById)
	dataSourceRouter.DELETE("/:id", sourceController.Delete)

	return service
}
