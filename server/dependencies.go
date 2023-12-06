package server

type DependencyManager struct {
}

type Initializable interface {
	Initialize(initializable ...Initializable)
}
