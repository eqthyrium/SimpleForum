package usecase

import (
	"SimpleForum/internal/domain"
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
)

// ToDo for SignUp:
// 1. Checking whether the entered input data is correct
// 2. Whether input email exist in the db
// 3. Hash the password
// 4. Insert into db (email, nickname, password(hashed), role(by default user))

func (app *Application) SignUp(nickname, email, password, oauth string) error {
	var user *entity.User = &entity.User{}

	if oauth == "direct" {
		//checking whether the input data is correct
		correctnessData := app.isItCorrect(nickname, email, password)
		if !correctnessData {
			return logger.ErrorWrapper("UseCase", "SignUp", "There is an invalid entered credentials of the client to be signed up", domain.ErrInvalidCredential)
		}
		nickname, email = makeItLower(nickname, email)

		// check is there exist such email
		recievedUser, err := app.ServiceDB.GetUserByEmail(email)
		if errors.Is(err, domain.ErrUserNotFound) {

			user.Nickname = nickname
			user.Email = email
			user.Password, err = hashPassword(password)
			if err != nil {
				return logger.ErrorWrapper("UseCase", "SignUp", "Failed to hash password", err)
			}
			user.Role = "User"

		} else {
			if err == nil {
				if recievedUser.Password != "" {
					return logger.ErrorWrapper("UseCase", "SignUp", "The client entered such credential which is already in the data base", domain.ErrInvalidCredential)
				}

				recievedUser.Password, err = hashPassword(password) // Here we have to write the Update sql query into db
				if err != nil {
					return logger.ErrorWrapper("UseCase", "SignUp", "Failed to hash password", err)
				}

				err := app.ServiceDB.UpdateUserPassword(recievedUser)
				if err != nil {
					return logger.ErrorWrapper("UseCase", "SignUp", "Failed to update user password", err)
				}

				return nil

			} else {
				return logger.ErrorWrapper("UseCase", "SignUp", "Failed to find user by email", err)
			}
		}

	} else {

		_, err := app.ServiceDB.GetUserByEmail(email)
		if !errors.Is(err, domain.ErrUserNotFound) {
			if err != nil {
				return logger.ErrorWrapper("UseCase", "SignUp", "Failed to find user by email", err)
			}
			return logger.ErrorWrapper("UseCase", "SignUp", "The client entered such credential which is already in the data base", domain.ErrInvalidCredential)
		}

		user.Nickname = truncatedNickname(email)
		user.Email = email
		user.Password = ""
		user.Role = "User"

	}

	err := app.ServiceDB.CreateUser(user)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "SignUp", "Failed to create user", err)
	}

	return nil
}

func (app *Application) isItCorrect(nickname, email, password string) bool {

	answerNickname := nicknameCheck(nickname)
	answerEmail := emailCheck(email)
	answerPassword := passwordCheck(password)

	if !(answerEmail && answerPassword && answerNickname) {
		return false
	}

	return true

}

func nicknameCheck(nickname string) bool {
	nicknameRegex := `^[a-zA-Z0-9]([a-zA-Z0-9._-]{1,18}[a-zA-Z0-9])?$`
	re := regexp.MustCompile(nicknameRegex)
	return re.MatchString(nickname)
}
func emailCheck(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func passwordCheck(password string) bool {

	if len(password) < 8 || len(password) > 32 {
		return false
	}

	for i := 0; i < len(password); i++ {
		if !(password[i] >= 32 && password[i] <= 126) {
			return false
		}
	}

	return true
}

func makeItLower(nickname, email string) (string, string) {
	return strings.ToLower(nickname), strings.ToLower(email)
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func truncatedNickname(email string) string {
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			return email[:i]
		}
	}
	return ""
}
