package usecase

import (
	"SimpleForum/internal/domain"
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
)

func (app *Application) GetLatestCommentaries(postId int) ([]entity.Commentaries, error) {
	return nil, nil
}

func (app *Application) CreateCommentary(userId, postId int, content string) error {
	validation := checkContent(content)
	if !validation {
		return logger.ErrorWrapper("UseCase", "CreateCommentary", "The provided content is not correponding to the application logic", domain.ErrNotValidContent)
	}

	err := app.ServiceDB.CreateCommentary(userId, postId, content)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "CreateCommentary", "Failed to create commentary", err)
	}
	err = app.ServiceDB.InsertReaction(userId, postId, "post", "comment")
	if err != nil {
		return logger.ErrorWrapper("UseCase", "CreateCommentary", "Failed to leave the commentaried reaction to particular post in the reaction db", err)
	}
	return nil
}

func checkContent(content string) bool {
	var existence bool
	for i := 0; i < len(content); i++ {
		if content[i] >= 33 && content[i] <= 126 {
			existence = true
		}
	}
	if !existence {
		return false
	}

	return true
}

func (app *Application) GetCertainPostsCommentaries(posts int) ([]entity.Commentaries, error) {
	comment, err := app.ServiceDB.GetCertainPostsCommentaries(posts)
	if err != nil {
		return nil, logger.ErrorWrapper("UseCase", "GetComments", "Failed to create commentary", err)
	}
	return comment, nil
}

func (app *Application) GetComments(userId int) ([]entity.Commentaries, error) {
	comment, err := app.ServiceDB.GetComments(userId)
	if err != nil {
		return nil, logger.ErrorWrapper("UseCase", "GetComments", "Failed to create commentary", err)
	}
	return comment, nil
}
