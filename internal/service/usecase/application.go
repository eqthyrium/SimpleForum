package usecase

import (
	"SimpleForum/internal/service/repository"
	"log"
)

type Application struct {
	ServiceDB *repository.ServiceRepository
	ErrorLog  *log.Logger
	InfoLog   *log.Logger
}
