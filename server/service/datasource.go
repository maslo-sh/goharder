package service

import (
	"proxy-engineering-thesis/model"
	"proxy-engineering-thesis/server/repository"
)

type DataSourceService interface {
	GetAll() ([]model.DataSource, error)
	GetById(id string) (*model.DataSource, error)
	Create(req *model.DataSource)
	Delete(id string) error
	Update(req *model.DataSource) error
}

type DataSourceServiceImpl struct {
	dataSourceRepository repository.DataSourceRepository
}

func NewDataSourceService(dataSourceRepository repository.DataSourceRepository) *DataSourceServiceImpl {
	return &DataSourceServiceImpl{dataSourceRepository: dataSourceRepository}
}

func (ps *DataSourceServiceImpl) GetById(id string) (*model.DataSource, error) {
	data, err := ps.dataSourceRepository.Get(id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ps *DataSourceServiceImpl) GetAll() ([]model.DataSource, error) {
	data, err := ps.dataSourceRepository.GetAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ps *DataSourceServiceImpl) Create(req *model.DataSource) {
	ps.dataSourceRepository.Create(req)
}

func (ps *DataSourceServiceImpl) Delete(id string) error {
	err := ps.dataSourceRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (ps *DataSourceServiceImpl) Update(dto *model.DataSource) error {
	return nil
}
