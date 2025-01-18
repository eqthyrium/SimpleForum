package usecase

import (
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
)

func (app *Application) GetAllCategories() ([]entity.Categories, error) {

	categories, err := app.ServiceDB.GetAllCategories()
	if err == nil {
		return categories, nil
	} else {
		return nil, logger.ErrorWrapper("UseCase", "GetAllCategories", "There is a problem getting all categories in db", err)
	}
}
