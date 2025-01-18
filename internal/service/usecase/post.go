package usecase

import (
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
)

func (app *Application) GetLatestPosts(requestedCategories []string) ([]entity.Posts, error) {

	posts, err := app.ServiceDB.GetLatestAllPosts(requestedCategories)
	if err != nil {
		return nil, logger.ErrorWrapper("UseCase", "GetLatestPosts", "There is problem with getting all the recent posts from the db", err)
	}

	return posts, nil
}
