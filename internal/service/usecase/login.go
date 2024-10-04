package usecase

import (
	"SimpleForum/internal/domain"
	"fmt"
)

// ToDo
// 1. Check whether the client exists
// 2. Making token to that client

func (app *Application) LogIn(email string, password string) error {

	receivedUser, err := app.ServiceDB.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("usecase-LogIn, %w", err)
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("usecase-LogIn, %w", err)
	}
	if hashedPassword != receivedUser.Password {
		return fmt.Errorf("usecase-LogIn, %w", domain.ErrUserNotFound)
	}

	token :=

	return nil
}
