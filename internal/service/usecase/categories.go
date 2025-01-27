package usecase

import (
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
	"strconv"
)

func (app *Application) GetAllCategories() ([]entity.Categories, error) {

	categories, err := app.ServiceDB.GetAllCategories()
	if err == nil {
		return categories, nil
	} else {
		return nil, logger.ErrorWrapper("UseCase", "GetAllCategories", "There is a problem getting all categories in db", err)
	}
}
func (app *Application) AddCategory(categoryName string) ([]entity.Categories, error) {

	categories, err := app.ServiceDB.AddCategory(categoryName)
	if err != nil {

		return nil, logger.ErrorWrapper("UseCase", "AddCategory", "There is a problem getting all categories in db", err)
	}
	return categories, nil

}
func (app *Application) DeleteCategory(categoryId string) ([]entity.Categories, error) {
	id, err := strconv.Atoi(categoryId)
	if err != nil {
		return nil, err
	}

	categories, err := app.ServiceDB.DeleteCategory(id)
	if err != nil {
		return nil, logger.ErrorWrapper("UseCase", "DeleteCategory", "There is a problem getting all categories in db", err)
	}

	return categories, nil

}
