package usecase

import (
	"SimpleForum/internal/domain/module"
)

type Application struct {
	ServiceDB module.DbModule
}

func NewUseCase(repoObject module.DbModule) *Application {
	return &Application{
		ServiceDB: repoObject,
	}
}
