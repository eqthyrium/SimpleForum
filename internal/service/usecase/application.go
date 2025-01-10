package usecase

import (
	"SimpleForum/internal/domain/module"
)

type UsecaseRepo struct {
	ServiceDB module.DbModule
}

func NewUseCase(repoObject module.DbModule) *UsecaseRepo {
	return &UsecaseRepo{
		ServiceDB: repoObject,
	}
}
