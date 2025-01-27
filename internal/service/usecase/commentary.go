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

func (app *Application) GetCertainCommentaryInfo(commentId int) (*entity.Commentaries, error) {
	commentary, err := app.ServiceDB.GetCertainCommentaryInfo(commentId)
	if err != nil {
		return nil, logger.ErrorWrapper("UseCase", "GetCertainCommentaryInfo", "Failed to retrieve certain comment info", err)
	}
	return commentary, nil
}

func (app *Application) DeleteCertainCommentary(userId, commentId int, role string) error {
	validation, err := app.ServiceDB.ValidateOfExistenceCertainCommentary(userId, commentId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "DeleteCertainCommentary", "There is problem with validating certain commentary", err)
	}
	if !validation && role == "User" {
		return logger.ErrorWrapper("UseCase", "DeleteCertainCommentary", "There is no such commentary of that user", domain.ErrCommentaryNotFound)
	}

	err = app.ServiceDB.DeleteCertainCommentary(commentId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "DeleteCertainCommentary", "There is problem with deleting certain commentary", err)
	}

	return nil

}

func (app *Application) EditCertainCommentary(userId, commentId int, content string) error {

	validation, err := app.ServiceDB.ValidateOfExistenceCertainCommentary(userId, commentId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "EditCertainCommentary", "There is problem with validating certain comment", err)
	}
	if !validation {
		return logger.ErrorWrapper("UseCase", "EditCertainCommentary", "There is no such comment of that user", domain.ErrCommentaryNotFound)
	}
	validation = checkContent(content)
	if !validation {
		return logger.ErrorWrapper("UseCase", "EditCertainCommentary", "There is no content", domain.ErrNotValidContent)
	}

	err = app.ServiceDB.UpdateEditedCommentary(userId, commentId, content)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "EditedCertainCommentary", "Failed to update the content of the commentary", err)
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

// ---------------------
func (app *Application) GetComments(userId int) ([]entity.Commentaries, error) {
	comment, err := app.ServiceDB.GetComments(userId)
	if err != nil {
		return nil, logger.ErrorWrapper("UseCase", "GetComments", "Failed to create commentary", err)
	}
	return comment, nil
}
