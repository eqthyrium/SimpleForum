package sqllite

import (
	"SimpleForum/internal/domain"
	"SimpleForum/pkg/logger"
	"database/sql"
	"errors"
)

func (rp *Repository) CreateUser(user *domain.User) error {
	statement := `INSERT INTO Users (Nickname, MemberIdentity, Password, Role) VALUES(?,?,?,?)`
	_, err := rp.DB.Exec(statement, user.Nickname, user.MemberIdentity, user.Password, user.Role)
	if err != nil {
		return logger.ErrorWrapper("Repository", "CreateUser", "The problem within the process of creation of the user in db", err)
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

func (rp *Repository) GetUserByEmail(memberIdentity string) (*domain.User, error) {

	statement := "SELECT UserId,MemberIdentity,Password, Role FROM Users WHERE MemberIdentity = ?"

	row := rp.DB.QueryRow(statement, memberIdentity)

	user := &domain.User{}

	err := row.Scan(&user.UserId, &user.MemberIdentity, &user.Password, &user.Role)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, logger.ErrorWrapper("Repository", "GetUserByEmail", "There is no such  user in the db", domain.ErrUserNotFound)
		} else {
			return nil, logger.ErrorWrapper("Repository", "GetUserByEmail", "The problem within the process of getting of the user by its MemberIdentity(email) in db", err)
		}
	}

	return user, nil
}
