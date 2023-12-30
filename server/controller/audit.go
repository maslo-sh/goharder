package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"proxy-engineering-thesis/model"
	"proxy-engineering-thesis/server/service"
)

type AuditController struct {
	auditService service.AuditService
}

func NewAuditController(auditService service.AuditService) *AuditController {
	return &AuditController{auditService: auditService}
}

func (ac *AuditController) PerformAudit(ctx *gin.Context) {
	var config model.AuditConfiguration
	id := ctx.Param("id")
	err := ctx.ShouldBindJSON(&config)
	if err != nil {
		fmt.Printf("failed to parse audit configuration object: %v\n", err)
		ctx.JSON(400, nil)
		return
	}
	auditResult, err := ac.auditService.PerformAudit(id, config)
	if err != nil {
		fmt.Printf("failed to perform audit: %v\n", err)
		ctx.JSON(400, nil)
		return
	}

	ctx.JSON(200, auditResult)
}
