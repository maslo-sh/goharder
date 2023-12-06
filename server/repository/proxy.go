package repository

import (
	"proxy-engineering-thesis/model"
)

type ProxyRepository interface {
	Create(req *model.ProxyDto) error
	Delete(id string) error
	Get(id string) (*model.ProxyDto, error)
	GetAll() ([]model.ProxyDto, error)
}

type ProxyRepositoryImpl struct {
	*DbContext
}

func NewProxyRepositoryImpl(dbCtx *DbContext) *ProxyRepositoryImpl {
	return &ProxyRepositoryImpl{dbCtx}
}

func (pr *ProxyRepositoryImpl) Create(req *model.ProxyDto) error {
	tx := pr.Db.Create(req)
	if err := tx.Error; err != nil {
		return tx.Error
	}
	return nil
}

func (pr *ProxyRepositoryImpl) Delete(id string) error {

	tx := pr.Db.Delete(&model.ProxyDto{}, id)
	if err := tx.Error; err != nil {
		return err
	}

	return nil
}

func (pr *ProxyRepositoryImpl) Get(id string) (*model.ProxyDto, error) {
	var proxyDto model.ProxyDto
	tx := pr.Db.First(&proxyDto, id)
	if err := tx.Error; err != nil {
		return nil, tx.Error
	}

	return &proxyDto, nil
}

func (pr *ProxyRepositoryImpl) GetAll() ([]model.ProxyDto, error) {
	var proxies []model.ProxyDto
	tx := pr.Db.Find(&proxies)
	if err := tx.Error; err != nil {
		return nil, tx.Error
	}

	return proxies, nil
}
