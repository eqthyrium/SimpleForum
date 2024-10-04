package sqllite

import (
	"SimpleForum/internal/domain"
	"database/sql"
	"errors"
	"fmt"
)

func (rp *Repository) CreateUser(user *domain.User) error {
	statement := `INSERT INTO User (Nickname, Email, Password, Role) VALUES(?,?,?,?)`
	_, err := rp.DB.Exec(statement, user.Nickname, user.Email, user.Password, user.Role)
	if err != nil {
		return fmt.Errorf("Repository-CreateUser: %w", err)
	}
	return nil
}

func (rp *Repository) UpdateUser(user *domain.User) error {
	return nil
}

func (rp *Repository) DeleteUser(user *domain.User) error {
	return nil
}

func (rp *Repository) GetUserByID(userId int) (domain.User, error) {
	return domain.User{}, nil
}

//func (rp *Repository) CheckUserByEmail(email string) (bool, error) {
//
//	statement := "SELECT Email FROM users WHERE email = ?"
//
//	row := rp.DB.QueryRow(statement, email)
//
//	user := &struct{ email string }{email: ""}
//
//	err := row.Scan(&user.email)
//	// Think about error
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return false, nil
//		} else {
//			return false, err
//		}
//	}
//	return true, nil
//}

func (rp *Repository) GetUserByEmail(email string) (*domain.User, error) {

	statement := "SELECT User_Id,Email,Password, Role FROM Users WHERE Email = ?"

	row := rp.DB.QueryRow(statement, email)

	user := &domain.User{}

	err := row.Scan(&user.UserId, &user.Email, &user.Password, &user.Role)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Repository-GetUserByEmail: %w", domain.ErrUserNotFound)
		} else {
			return nil, fmt.Errorf("Repository-GetUserByEmail: %w", err)
		}
	}

	return user, nil
}
