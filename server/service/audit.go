package service

import (
	"fmt"
	"proxy-engineering-thesis/internal/audit/sql"
	"proxy-engineering-thesis/model"
)

type AuditService interface {
	PerformAudit(string, model.AuditConfiguration) (model.AuditData, error)
}

type AuditServiceImpl struct {
	dataSourceService DataSourceService
}

func NewAuditService(dsService DataSourceService) *AuditServiceImpl {
	return &AuditServiceImpl{dataSourceService: dsService}
}

func (as AuditServiceImpl) PerformAudit(id string, config model.AuditConfiguration) (model.AuditData, error) {
	ds, err := as.dataSourceService.GetById(id)
	if err != nil {
		fmt.Printf("failed to retrieve auditted datasource data: %v\n", err)
	}

	dsConnData := sql.DataSourceConnectionData{
		Host:     ds.Hostname,
		Port:     ds.Port,
		Username: ds.Username,
		Password: ds.Password,
	}

	auditResult, err := sql.PerformAudit(dsConnData, config)
	if err != nil {
		fmt.Printf("failed to perform audit: %v\n", err)
	}
	return auditResult, nil
}
