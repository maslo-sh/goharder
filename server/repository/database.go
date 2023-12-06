package repository

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"proxy-engineering-thesis/internal/utils"
	"proxy-engineering-thesis/model"
)

var tables = map[string]interface{}{
	"proxies":     &model.ProxyDto{},
	"datasources": &model.DataSource{},
}

type DbContext struct {
	Db *gorm.DB
}

func NewDbContext() *DbContext {
	return &DbContext{connectToDatabase()}
}

func (db *DbContext) SetUpSchema() {
	db.Db.AutoMigrate(&model.ProxyDto{}, &model.DataSource{})
}

func connectToDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	utils.ErrorPanic(err)
	return db
}
