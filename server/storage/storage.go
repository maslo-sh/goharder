package storage

import (
	"fmt"
	"proxy-engineering-thesis/internal/proxy/sql"
)

type ProxiesStorage struct {
	Proxies map[string]*sql.ProxyConfiguration
}

func NewProxiesStorage() *ProxiesStorage {
	return &ProxiesStorage{
		make(map[string]*sql.ProxyConfiguration),
	}
}

func (p *ProxiesStorage) AddProxyToStorage(proxy *sql.ProxyConfiguration, id string) string {
	p.Proxies[id] = proxy
	return id
}

func (p *ProxiesStorage) RemoveProxyFromStorage(id string) error {
	proxyConf := p.Proxies[id]
	if proxyConf == nil {
		return fmt.Errorf("no proxy instance with given id: %s", id)
	}
	proxyConf.CloseSessions()
	proxyConf.Listener.Close()
	delete(p.Proxies, id)
	return nil
}
