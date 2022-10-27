package server

type Object interface {
	GetType() string
	GetName() string
}
