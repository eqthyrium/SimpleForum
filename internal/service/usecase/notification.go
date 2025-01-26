package usecase

import (
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
)

func (app *Application) GetNotifications(userId int) ([]entity.Notifications, error) {

	notifications, err := app.ServiceDB.GetReactionsNotificationOfPosts(userId)
	if err != nil {
		return nil, logger.ErrorWrapper("UseCase", "GetNotifications", "Failed to get all recent notifications", err)
	}

	for i := 0; i < len(notifications); i++ {
		if notifications[i].Action == "L" {
			notifications[i].Action = "liked"
		} else if notifications[i].Action == "D" {
			notifications[i].Action = "disliked"
		} else if notifications[i].Action == "C" {
			notifications[i].Action = "commented"
		}
	}

	return notifications, nil
}
