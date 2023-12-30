package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"net/http"
	"proxy-engineering-thesis/internal/utils"
	"proxy-engineering-thesis/server/controller"
	"proxy-engineering-thesis/server/repository"
	"proxy-engineering-thesis/server/service"
	"proxy-engineering-thesis/server/storage"
	"time"
)

func StartServer() {
	container := dig.New()
	declareDependencies(container)

	dbContextInitialization := func(dbCtx *repository.DbContext) {
		dbCtx.SetUpSchema()
	}

	routeDeclaration := func(proxyController *controller.ProxyController, dsController *controller.DataSourceController, auditController *controller.AuditController) {
		routes := NewRouter(proxyController, dsController, auditController)

		server := &http.Server{
			Addr:           ":8888",
			Handler:        routes,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		err := server.ListenAndServe()
		utils.ErrorPanic(err)
	}

	if err := container.Invoke(dbContextInitialization); err != nil {
		panic(err)
	}

	if err := container.Invoke(routeDeclaration); err != nil {
		panic(err)
	}

}

func declareDependencies(container *dig.Container) {
	container.Provide(repository.NewDbContext)
	container.Provide(func(db *repository.DbContext) repository.ProxyRepository {
		return repository.NewProxyRepositoryImpl(db)
	})
	container.Provide(storage.NewProxiesStorage)
	container.Provide(func(db *repository.DbContext) repository.DataSourceRepository {
		return repository.NewDataSourceRepositoryImpl(db)
	})
	container.Provide(func(dsRepo repository.DataSourceRepository) service.DataSourceService {
		return service.NewDataSourceService(dsRepo)
	})
	container.Provide(func(proxyRepo repository.ProxyRepository, dsService service.DataSourceService, proxyStorage *storage.ProxiesStorage) service.ProxyService {
		return service.NewProxyService(proxyRepo, dsService, proxyStorage)
	})
	container.Provide(func(dsService service.DataSourceService) *controller.DataSourceController {
		return controller.NewDataSourceController(dsService)
	})
	container.Provide(func(dsService service.DataSourceService) service.AuditService {
		return service.NewAuditService(dsService)
	})
	container.Provide(func(proxyService service.ProxyService) *controller.ProxyController {
		return controller.NewProxyController(proxyService)
	})
	container.Provide(func(auditService service.AuditService) *controller.AuditController {
		return controller.NewAuditController(auditService)
	})
	container.Provide(func(proxyController *controller.ProxyController, dsController *controller.DataSourceController, auditController *controller.AuditController) *gin.Engine {
		return NewRouter(proxyController, dsController, auditController)
	})
}
