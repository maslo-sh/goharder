package model

import (
	"fmt"
)

type Address struct {
	Hostname string
	Port     string
}

func (a Address) CreateHostString() string {
	return fmt.Sprintf("%s:%s", a.Hostname, a.Port)
}
