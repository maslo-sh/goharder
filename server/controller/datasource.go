package controller

import (
	"github.com/gin-gonic/gin"
	"proxy-engineering-thesis/model"
	"proxy-engineering-thesis/server/service"
)

type DataSourceController struct {
	dataSourceService service.DataSourceService
}

func NewDataSourceController(dataSourceService service.DataSourceService) *DataSourceController {
	return &DataSourceController{dataSourceService: dataSourceService}
}

func (pc *DataSourceController) Create(ctx *gin.Context) {
	var req model.DataSource
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(400, nil)
	}
	pc.dataSourceService.Create(&req)

	ctx.JSON(200, req)
}

func (pc *DataSourceController) Delete(ctx *gin.Context) {
	var name string
	err := ctx.ShouldBindQuery(&name)
	if err != nil {
		ctx.JSON(400, nil)
	}
	pc.dataSourceService.Delete(name)
	ctx.JSON(200, nil)
}

func (pc *DataSourceController) GetAll(ctx *gin.Context) {
	//var proxies []dto.ProxyDto
}

func (pc *DataSourceController) FindById(ctx *gin.Context) {
	//var proxy dto.ProxyDto
}
