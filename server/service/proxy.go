package service

import (
	"github.com/google/uuid"
	"proxy-engineering-thesis/internal/proxy/sql"
	"proxy-engineering-thesis/model"
	"proxy-engineering-thesis/server/repository"
	"proxy-engineering-thesis/server/storage"
	"strconv"
)

type ProxyService interface {
	GetAll() ([]model.ProxyDto, error)
	GetById(id string) (*model.ProxyDto, error)
	Create(req *model.ProxyDto)
	Delete(id string) error
	Update(req *model.ProxyDto) error
	StartProxy(id string) (string, error)
	StopProxy(id string) error
}

type ProxyServiceImpl struct {
	proxyRepository   repository.ProxyRepository
	dataSourceService DataSourceService
	proxiesStorage    *storage.ProxiesStorage
}

func NewProxyService(
	proxyRepository repository.ProxyRepository,
	sourceService DataSourceService,
	storage *storage.ProxiesStorage) *ProxyServiceImpl {
	return &ProxyServiceImpl{proxyRepository: proxyRepository, dataSourceService: sourceService, proxiesStorage: storage}
}

func (ps *ProxyServiceImpl) GetById(id string) (*model.ProxyDto, error) {
	data, err := ps.proxyRepository.Get(id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ps *ProxyServiceImpl) GetAll() ([]model.ProxyDto, error) {
	data, err := ps.proxyRepository.GetAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ps *ProxyServiceImpl) Create(req *model.ProxyDto) {
	ps.proxyRepository.Create(req)
}

func (ps *ProxyServiceImpl) Delete(id string) error {
	err := ps.proxyRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProxyServiceImpl) Update(dto *model.ProxyDto) error {
	return nil
}

func (ps *ProxyServiceImpl) StartProxy(id string) (string, error) {
	proxyDto, err := ps.proxyRepository.Get(id)
	if err != nil {
		return "", err
	}

	dsId := strconv.FormatUint(uint64(proxyDto.DataSourceID), 10)

	ds, err := ps.dataSourceService.GetById(dsId)

	if err != nil {
		return "", err
	}

	proxyConfig := sql.NewProxy(*proxyDto, *ds)
	processId := uuid.New()
	ps.proxiesStorage.AddProxyToStorage(proxyConfig, processId.String())
	go proxyConfig.Start()
	return processId.String(), nil
}

func (ps *ProxyServiceImpl) StopProxy(id string) error {
	err := ps.proxiesStorage.RemoveProxyFromStorage(id)
	if err != nil {
		return err
	}
	return nil
}
