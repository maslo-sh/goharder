package datasource

import (
	"fmt"
	"proxy-engineering-thesis/model"
)

type SqlDataSource struct {
	model.Address
	User     string
	Password string
	DbName   string
}

func (ds SqlDataSource) GetTargetAddress() string {
	//return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", ds.User, ds.Password, ds.Hostname, ds.Port, ds.DbName)
	return fmt.Sprintf("%s:%s", ds.Hostname, ds.Port)
}

func (ds SqlDataSource) Scan() {

}

func NewSqlDataSource(host string, port string, user string, password string, dbName string) SqlDataSource {
	return SqlDataSource{
		Address: model.Address{
			Hostname: host,
			Port:     port,
		},
		User:     user,
		Password: password,
		DbName:   dbName,
	}
}
