package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"proxy-engineering-thesis/model"
	"proxy-engineering-thesis/server/service"
	"strconv"
)

type ProxyController struct {
	proxyService service.ProxyService
}

func NewProxyController(proxyService service.ProxyService) *ProxyController {
	return &ProxyController{proxyService: proxyService}
}

func (pc *ProxyController) Create(ctx *gin.Context) {
	var req model.ProxyDto
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
	}
	pc.proxyService.Create(&req)

	ctx.JSON(200, req)
}

func (pc *ProxyController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	pc.proxyService.Delete(id)
	ctx.JSON(200, nil)
}

func (pc *ProxyController) GetAll(ctx *gin.Context) {
	proxies, err := pc.proxyService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "")
	}

	ctx.JSON(200, proxies)
}

func (pc *ProxyController) FindById(ctx *gin.Context) {
	id := ctx.Param("id")
	proxy, err := pc.proxyService.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "")
	}

	ctx.JSON(200, proxy)
}

func (pc *ProxyController) StartProxy(ctx *gin.Context) {
	id := ctx.Param("id")
	proxyId, err := pc.proxyService.StartProxy(id)
	if err != nil {
		ctx.Data(http.StatusBadRequest, "text/plain; charset=utf-8", []byte("No configuration found"))
	}

	ctx.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(proxyId))
}

func (pc *ProxyController) StopProxy(ctx *gin.Context) {
	id := ctx.Param("id")
	err := pc.proxyService.StopProxy(id)
	if err != nil {
		ctx.Data(http.StatusBadRequest, "text/plain; charset=utf-8", []byte(err.Error()))
	}
}

func (pc *ProxyController) GetProxySessionsCount(ctx *gin.Context) {
	id := ctx.Param("id")
	count, err := pc.proxyService.GetProxySessionsCount(id)
	if err != nil {
		ctx.Data(http.StatusBadRequest, "text/plain; charset=utf-8", []byte(err.Error()))
	}

	ctx.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(strconv.Itoa(count)))
}
