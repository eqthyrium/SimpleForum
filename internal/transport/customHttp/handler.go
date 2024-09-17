package customHttp

import "log"

type httpModule interface {
	// Here we write what kind of services can be used in the http handler
}

type Handler struct {
	Service  httpModule
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

func NewTransportHttpHandler(service httpModule) *Handler {
	return &Handler{Service: service}
}
