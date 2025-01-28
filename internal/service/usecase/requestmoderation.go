package usecase

import (
	"SimpleForum/internal/domain"
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
)

func (app *Application) GetAllRequests() ([]entity.ReportInfo, error) {
	requests, err := app.ServiceDB.GetAllRequests()
	if err != nil {
		return nil, logger.ErrorWrapper("UseCase", "GetAllRequests", "Failed to get all requests from database", err)
	}
	return requests, nil
}

func (app *Application) RequestToBeModerator(userId int) error {

	validation, err := app.ServiceDB.IsItRequestedToBeModerator(userId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "RequestToBeModerator", "Failed to check the existence of a request to be moderator for a user", err)
	}
	if validation {
		return logger.ErrorWrapper("UseCase", "RequestToBeModerator", "There is already such request to be moderator for a user in db", domain.ErrRepeatedRequest)
	}

	err = app.ServiceDB.CreateRequestToBeModerator(userId, -1)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "RequestToBeModerator", "Failed to set of a request to be moderator for a user", err)
	}

	return nil

}

func (app *Application) AcceptRequestToBeModerator(userId int) error {
	role, err := app.ServiceDB.GetUsersRole(userId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "AcceptRequestToBeModerator", "Failed to get users role", err)
	}
	if role != "User" {
		return logger.ErrorWrapper("UseCase", "AcceptRequestToBeModerator", "The client sent request to be moderator cannot be done, because it's role in db is not User", domain.ErrInvalidOperation)
	}

	err = app.ServiceDB.UpdateRoleOfUser(userId, "Moderator")
	if err != nil {
		return logger.ErrorWrapper("UseCase", "AcceptRequestToBeModerator", "Failed to set of a request to be moderator for a user", err)
	}

	err = app.ServiceDB.DeleteRequestToBeModerator(userId, -1)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "AcceptRequestToBeModerator", "Failed to delete request to moderator for a user", err)
	}

	return nil

}

func (app *Application) DeclineRequestToBeModerator(userId int) error {
	role, err := app.ServiceDB.GetUsersRole(userId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "DeclineRequestToBeModerator", "Failed to get users role", err)
	}
	if role != "User" {
		return logger.ErrorWrapper("UseCase", "DeclineRequestToBeModerator", "The client sent request to be moderator cannot be done, because it's role in db is not User", domain.ErrInvalidOperation)
	}
	err = app.ServiceDB.DeleteRequestToBeModerator(userId, -1)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "DeclineRequestToBeModerator", "Failed to delete request to moderator for a user", err)
	}

	return nil
}
