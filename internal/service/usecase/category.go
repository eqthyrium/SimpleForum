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

func (app *Application) CreateCategory(categoryName string) error {

	err := app.ServiceDB.CreateCategory(categoryName)
	if err != nil {

		return logger.ErrorWrapper("UseCase", "AddCategory", "There is a problem getting all categories in db", err)
	}
	return nil

}

// Before deleting a particular category, i have to check out whether there is only one post which is related to it, if so i have to delete that post too.
func (app *Application) DeleteCategory(categoryId int) error {

	postIds, err := app.ServiceDB.GetPostIDsOfCertainCategory(categoryId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "DeleteCategory", "There is a problem getting all postcategory object by certain category", err)
	}

	for i := 0; i < len(postIds); i++ {
		categoryIdsOfCertainPost, err := app.ServiceDB.GetCategoriesOfCertainPost(postIds[i])
		if len(categoryIdsOfCertainPost) == 1 && categoryIdsOfCertainPost[0] == categoryId {
			err = app.ServiceDB.DeleteCertainPost(postIds[i])
			if err != nil {
				return logger.ErrorWrapper("UseCase", "DeleteCategory", "There is a problem deleting post in db", err)
			}

		}
	}

	err = app.ServiceDB.DeleteCategory(categoryId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "DeleteCategory", "There is a problem getting all categories in db", err)
	}

	return nil

}
