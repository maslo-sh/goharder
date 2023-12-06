package datasource

import "proxy-engineering-thesis/model"

type LdapDataSource struct {
	model.Address
}

func (ds LdapDataSource) GetTargetAddress() string {
	return "ldap://" + ds.Hostname + ":" + ds.Port
}

func (ds LdapDataSource) Scan() {

}
