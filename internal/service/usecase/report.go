package usecase

import (
	"SimpleForum/internal/domain"
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
)

func (app *Application) GetAllReports() ([]entity.ReportInfo, error) {
	reports, err := app.ServiceDB.GetAllReports()
	if err != nil {
		return nil, logger.ErrorWrapper("UseCase", "GetAllReports", "Failed to get all reports from database", err)
	}
	return reports, nil
}

func (app *Application) ReportPost(userId, postId int) error {
	validation, err := app.ServiceDB.IsItReportedPost(userId, postId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "ReportPost", "Failed to check the existence of a report to admin", err)
	}
	if validation {
		return logger.ErrorWrapper("UseCase", "ReportPost", "There is already such report from moderator to admin in db", domain.ErrRepeatedRequest)
	}

	err = app.ServiceDB.CreateReport(userId, postId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "ReportPost", "Failed to create report of the moderator to admin", err)
	}

	return nil
}

func (app *Application) AcceptCertainReport(userId, postId int) error {
	role, err := app.ServiceDB.GetUsersRole(userId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "AcceptCertainReport", "Failed to get role for user", err)
	}
	if role != "Moderator" {
		return logger.ErrorWrapper("UseCase", "AcceptCertainReport", "The non-moderator cannot report on post", domain.ErrInvalidOperation)
	}

	err = app.ServiceDB.DeleteCertainReport(userId, postId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "AcceptCertainReport", "Failed to delete report of the moderator", err)
	}

	err = app.ServiceDB.DeleteCertainPost(postId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "AcceptCertainReport", "Failed to delete report of the post", err)
	}
	return nil
}

func (app *Application) DeclineCertainReport(userId, postId int) error {
	role, err := app.ServiceDB.GetUsersRole(userId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "AcceptCertainReport", "Failed to get role for user", err)
	}
	if role != "Moderator" {
		return logger.ErrorWrapper("UseCase", "AcceptCertainReport", "The non-moderator cannot report on post", domain.ErrInvalidOperation)
	}
	err = app.ServiceDB.DeleteCertainReport(userId, postId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "AcceptCertainReport", "Failed to delete report of the moderator", err)
	}
	return nil
}
