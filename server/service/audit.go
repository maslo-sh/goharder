package service

import (
	"fmt"
	"proxy-engineering-thesis/internal/audit/relational"
	relationalModel "proxy-engineering-thesis/model/relational"
)

type AuditService interface {
	PerformAudit(string, relationalModel.AuditConfiguration) (relationalModel.AuditData, error)
}

type AuditServiceImpl struct {
	dataSourceService DataSourceService
}

func NewAuditService(dsService DataSourceService) *AuditServiceImpl {
	return &AuditServiceImpl{dataSourceService: dsService}
}

func (as AuditServiceImpl) PerformAudit(id string, config relationalModel.AuditConfiguration) (relationalModel.AuditData, error) {
	ds, err := as.dataSourceService.GetById(id)
	if err != nil {
		fmt.Printf("failed to retrieve auditted datasource data: %v\n", err)
	}

	dsConnData := relational.DataSourceConnectionData{
		Host:     ds.Hostname,
		Port:     ds.Port,
		Username: ds.Username,
		Password: ds.Password,
	}

	auditResult, err := relational.PerformAudit(dsConnData, config)
	if err != nil {
		fmt.Printf("failed to perform audit: %v\n", err)
	}
	return auditResult, nil
}
