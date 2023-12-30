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
	err = pc.dataSourceService.Delete(name)
	if err != nil {
		ctx.JSON(400, nil)
	}
	ctx.JSON(200, nil)
}

func (pc *DataSourceController) GetAll(ctx *gin.Context) {
	dataSources, err := pc.dataSourceService.GetAll()
	if err != nil {
		ctx.JSON(400, nil)
	}

	ctx.JSON(200, dataSources)
}

func (pc *DataSourceController) FindById(ctx *gin.Context) {
	//var proxy dto.ProxyDto
}
