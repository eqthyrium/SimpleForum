package usecase

import (
	"SimpleForum/internal/service/repository"
)

type Application struct {
	ServiceDB *repository.ServiceRepository
}
