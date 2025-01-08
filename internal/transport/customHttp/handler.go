package customHttp

import (
	"SimpleForum/internal/domain/module"
	"SimpleForum/pkg/logger"
)

type HandlerHttp struct {
	Service module.HttpModule
}

var customLogger = logger.NewLogger().GetLoggerObject("../logging/info.log", "../logging/error.log", "../logging/debug.log")

func NewTransportHttpHandler(ServiceObject module.HttpModule) *HandlerHttp {
	handlerObject := &HandlerHttp{
		Service: ServiceObject,
	}

	return handlerObject
}
