package sqllite

import (
	"SimpleForum/internal/domain"
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
	"database/sql"
	"errors"
)

// MemberIdentity is the Email!!!!!!!

func (rp *Repository) CreateUser(user *entity.User) error {
<<<<<<< HEAD
	statement := `INSERT INTO Users (Nickname, Email, Password, Role) VALUES(?,?,?,?)`
=======
	statement := `
	INSERT INTO "user" (
		nickname, 
		email, 
		password, 
		role
	) 
	VALUES(?,?,?,?)
	`
>>>>>>> 4b447443e2c202f7fd81126d1feb345e0481d725
	_, err := rp.DB.Exec(statement, user.Nickname, user.Email, user.Password, user.Role)
	if err != nil {
		return logger.ErrorWrapper("Repository", "CreateUser", "The problem within the process of creation of the user in db", err)
	}
	return nil
}

<<<<<<< HEAD
func (rp *Repository) UpdateUserPassword(user *entity.User) error {
	statement := `UPDATE Users SET Password = ? WHERE Email = ?`
	_, err := rp.DB.Exec(statement, user.Password, user.Email)
	if err != nil {
		return logger.ErrorWrapper("Repository", "UpdateUserPassword", "The problem within the process of updating the password of the user in db", err)
=======
func (rp *Repository) UpdateUser(user *entity.User) error {
	statement := `
	UPDATE "user"
	SET role = ?
	WHERE user_id = ?
	`
	_, err := rp.DB.Exec(statement, user.Role, user.UserId)
	if err != nil {
		return logger.ErrorWrapper("Repository", "UpdateUser", "The problem in the process of Update user", err)
>>>>>>> 4b447443e2c202f7fd81126d1feb345e0481d725
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

func (rp *Repository) GetUserByEmail(email string) (*entity.User, error) {

	statement := "SELECT UserId,Email,Password, Role FROM Users WHERE Email = ?"

	row := rp.DB.QueryRow(statement, email)

	user := &entity.User{}

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
