package model

import "gorm.io/gorm"

type DataSource struct {
	gorm.Model
	Address    `gorm:"embedded"`
	Credential `gorm:"embedded"`
	Proxies    []ProxyDto
}
