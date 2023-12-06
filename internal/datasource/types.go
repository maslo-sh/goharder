package datasource

type Connectable interface {
	GetTargetAddress() string
}

type Scanable interface {
	Scan()
}
