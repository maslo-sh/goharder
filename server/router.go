package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"proxy-engineering-thesis/server/controller"
)

const WebAppUrl = "http://localhost:3000"

func NewRouter(
	proxyController *controller.ProxyController,
	sourceController *controller.DataSourceController,
	auditController *controller.AuditController) *gin.Engine {
	service := gin.Default()

	service.NoRoute(WebAppReverseProxy)

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

	auditRouter := router.Group("/audit")
	auditRouter.POST("/:id", auditController.PerformAudit)

	return service
}

func WebAppReverseProxy(c *gin.Context) {
	remote, _ := url.Parse(WebAppUrl)
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL = c.Request.URL
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
