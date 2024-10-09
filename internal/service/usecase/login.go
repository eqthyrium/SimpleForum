package usecase

import (
	"SimpleForum/internal/domain"
	"SimpleForum/internal/service/auth"
	"fmt"
)

// ToDo
// 1. Check whether the client exists
// 2. Making token to that client

func (app *Application) LogIn(email string, password string) (string, string, error) {

	receivedUser, err := app.ServiceDB.GetUserByEmail(email) // The handler side must check whether its error is ErrUserNotFound error, in order to be adjusted in giving back webpage
	if err != nil {
		return "", "", fmt.Errorf("usecase-LogIn, %w", err)
	}
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return "", "", fmt.Errorf("usecase-LogIn, %w", err)
	}

	if hashedPassword != receivedUser.Password {
		return "", "", fmt.Errorf("usecase-LogIn, %w", domain.ErrUserNotFound)
	}

	tokenSignature, err := auth.CreateToken(receivedUser.UserId, receivedUser.Role)
	if err != nil {
		return "", "", fmt.Errorf("usecase-LogIn, %w", err)
	}

	return tokenSignature, receivedUser.Role, nil
}
