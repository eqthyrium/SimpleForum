package usecase

import (
	"SimpleForum/internal/domain"
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
	"errors"
)

func (app *Application) GetCertainUsers() ([]entity.Users, error) {
	users, err := app.ServiceDB.GetCertainUsers()
	if err != nil {
		return nil, logger.ErrorWrapper("UseCase", "GetCertainUsers", "Failed to get particular set of users", err)
	}
	return users, nil
}

func (app *Application) ChangeRole(userId int, action string) error {
	role, err := app.ServiceDB.GetUsersRole(userId)
	if errors.Is(err, domain.ErrUserNotFound) {
		return logger.ErrorWrapper("UseCase", "ChangeRole", "Failed to change  role, due to invalid operation", domain.ErrInvalidOperation)
	} else if err != nil {
		return logger.ErrorWrapper("UseCase", "ChangeRole", "Failed to change  role", err)
	}

	if action == "promote" && role == "User" {

		err = app.ServiceDB.UpdateRoleOfUser(userId, "Moderator")
		if err != nil {
			return logger.ErrorWrapper("UseCase", "ChangeRoleOfUser", "Failed to promote the user to moderator", err)
		}
		validation, err := app.ServiceDB.CheckExistenceOfSuchReport(userId, -1)
		if err != nil {
			return logger.ErrorWrapper("UseCase", "ChangeRoleOfUser", "Failed to check the report existence", err)
		}
		if validation {
			err = app.ServiceDB.DeleteRequestToBeModerator(userId, -1)
			if err != nil {
				return logger.ErrorWrapper("UseCase", "ChangeRoleOfUser", "Failed to delete the request to moderator", err)
			}
		}

	} else if action == "demote" && role == "Moderator" {

		err := app.ServiceDB.UpdateRoleOfUser(userId, "User")
		if err != nil {
			return logger.ErrorWrapper("UseCase", "ChangeRoleOfUser", "Failed to demote the moderator to user", err)
		}

	} else {
		return logger.ErrorWrapper("UseCase", "ChangeRole", "Failed to change  role, due to invalid operation", domain.ErrInvalidOperation)
	}

	return nil
}
