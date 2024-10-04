package customHttp

import "log"

type httpModule interface {
	SignUp(nickname, email, password string) (err error)
	LogIn(email, password string) (err error)
	//Login(email, password string) (err error)
	// Here we write what kind of services can be used in the http handler
}

type Handler struct {
	Service  httpModule
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

func NewTransportHttpHandler(Service httpModule) httpModule {
	return Service
}
