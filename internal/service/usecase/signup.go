package usecase

import "SimpleForum/internal/domain"

// Email must be unique

func (app *Application) SignUp(user *domain.User) error {
	//  checking whether the input data is correct
	//  check is there exist such email, or nickname
	//

	return nil
}

func (app *Application) SignUpChecking() (bool, error) {

	return false, nil
}

func (app *Application) correctnessData() (map[string]bool, error) {

	// Email

	// Nickname

	// Password

	return nil, nil
}
