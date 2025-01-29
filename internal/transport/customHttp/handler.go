package customHttp

import (
	"SimpleForum/internal/domain/module"
	"SimpleForum/pkg/logger"
)

var customLogger = logger.NewLogger().GetLoggerObject("./logging/logger.log", "./logging/logger.log", "./logging/logger.log")

type HandlerHttp struct {
	Service module.HttpModule
}

func NewTransportHttpHandler(ServiceObject module.HttpModule) *HandlerHttp {
	handlerObject := &HandlerHttp{
		Service: ServiceObject,
	}

	return handlerObject
}
