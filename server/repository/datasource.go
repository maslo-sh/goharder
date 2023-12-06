package repository

import (
	"proxy-engineering-thesis/model"
)

type DataSourceRepository interface {
	Create(req *model.DataSource) error
	Delete(id string) error
	Get(id string) (*model.DataSource, error)
	GetAll() ([]model.DataSource, error)
}

type DataSourceRepositoryImpl struct {
	*DbContext
}

func NewDataSourceRepositoryImpl(dbCtx *DbContext) *DataSourceRepositoryImpl {
	return &DataSourceRepositoryImpl{dbCtx}
}

func (pr *DataSourceRepositoryImpl) Create(req *model.DataSource) error {
	tx := pr.Db.Create(req)
	if err := tx.Error; err != nil {
		return tx.Error
	}
	return nil
}

func (pr *DataSourceRepositoryImpl) Delete(id string) error {

	tx := pr.Db.Delete(&model.DataSource{}, id)
	if err := tx.Error; err != nil {
		return err
	}

	return nil
}

func (pr *DataSourceRepositoryImpl) Get(id string) (*model.DataSource, error) {
	var dsDto model.DataSource
	tx := pr.Db.First(&dsDto, id)
	if err := tx.Error; err != nil {
		return nil, tx.Error
	}

	return &dsDto, nil
}

func (pr *DataSourceRepositoryImpl) GetAll() ([]model.DataSource, error) {
	return nil, nil
}
