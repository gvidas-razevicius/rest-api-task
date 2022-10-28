package server

// Interface to be implemented by any object that is going to be handled by the server and client
type Object interface {
	GetType() string
	GetName() string
}
