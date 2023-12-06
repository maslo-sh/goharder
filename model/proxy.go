package model

import "gorm.io/gorm"

type ProxyDto struct {
	gorm.Model
	Name         string
	Address      `gorm:"embedded"`
	DataSourceID uint
}

func (ProxyDto) TableName() string {
	return "proxies"
}
