package sqllite

import (
	"SimpleForum/internal/domain"
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
	"database/sql"
	"errors"
)

// MemberIdentity is the Email!!!!!!!

func (rp *Repository) CreateUser(user *entity.Users) error {
	statement := `INSERT INTO Users (Nickname, Email, Password, Role) VALUES(?,?,?,?)`
	_, err := rp.DB.Exec(statement, user.Nickname, user.Email, user.Password, user.Role)
	if err != nil {
		return logger.ErrorWrapper("Repository", "CreateUser", "The problem within the process of creation of the user in db", err)
	}
	return nil
}

func (rp *Repository) UpdateUserPassword(user *entity.Users) error {
	statement := `UPDATE Users SET Password = ? WHERE Email = ?`
	_, err := rp.DB.Exec(statement, user.Password, user.Email)
	if err != nil {
		return logger.ErrorWrapper("Repository", "UpdateUserPassword", "The problem within the process of updating the password of the user in db", err)
	}
	return nil
}

func (rp *Repository) GetCertainUsers() ([]entity.Users, error) {
	statement :=
		`
		SELECT UserId,Email,  Role FROM Users WHERE Role = "User" OR Role = "Moderator";
	    `

	rows, err := rp.DB.Query(statement)
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetCertainUsers", "Failed to get certain type of users from db", err)
	}
	var users []entity.Users
	for rows.Next() {
		var user entity.Users
		err := rows.Scan(&user.UserId, &user.Email, &user.Role)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "GetCertainUsers", "Failed to scan through already got result in db", err)
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetCertainUsers", "There is an error of the row in db", err)
	}
	if err := rows.Close(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetCertainUsers", "Failed to close the row of db", err)
	}

	return users, nil
}

func (rp *Repository) GetUsersRole(userId int) (string, error) {
	statement := `
			SELECT Role FROM Users WHERE UserId = ?;
`
	row := rp.DB.QueryRow(statement, userId)
	var role string
	err := row.Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", domain.ErrUserNotFound
		}
		return "", logger.ErrorWrapper("Repository", "GetUsersRole", "Failed to get role of user", err)
	}
	return role, nil
}

func (rp *Repository) UpdateRoleOfUser(userId int, role string) error {
	statement := `UPDATE Users SET Role = ? WHERE UserId = ?`
	_, err := rp.DB.Exec(statement, role, userId)
	if err != nil {
		return logger.ErrorWrapper("Repository", "UpdateRoleOfUser", "The problem within the process of updating the role of the user in db", err)
	}
	return nil
}

//func (rp *Repository) UpdateUser(user *entity.User) error {
//	return nil
//}
//
//func (rp *Repository) DeleteUser(user *entity.User) error {
//	return nil
//}
//
//func (rp *Repository) GetUserByID(userId int) (entity.User, error) {
//	return entity.User{}, nil
//}
//
////func (rp *Repository) CheckUserByEmail(email string) (bool, error) {
////
////	statement := "SELECT Email FROM users WHERE email = ?"
////
////	row := rp.DB.QueryRow(statement, email)
////
////	user := &struct{ email string }{email: ""}
////
////	err := row.Scan(&user.email)
////	// Think about error
////	if err != nil {
////		if errors.Is(err, sql.ErrNoRows) {
////			return false, nil
////		} else {
////			return false, err
////		}
////	}
////	return true, nil
////}

func (rp *Repository) GetUserByEmail(email string) (*entity.Users, error) {

	statement := "SELECT UserId,Email,Password, Role FROM Users WHERE Email = ?"

	row := rp.DB.QueryRow(statement, email)

	user := &entity.Users{}

	err := row.Scan(&user.UserId, &user.Email, &user.Password, &user.Role)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, logger.ErrorWrapper("Repository", "GetUserByEmail", "There is no such  user in the db", domain.ErrUserNotFound)
		} else {
			return nil, logger.ErrorWrapper("Repository", "GetUserByEmail", "The problem within the process of getting of the user by its email in db", err)
		}
	}

	return user, nil
}
