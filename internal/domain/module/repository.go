package module

import (
	"SimpleForum/internal/domain/entity"
)

type DbModule interface {
	userRepository
	//PostRepository
	//CommentRepository
	//CategoryRepository
	//PostCategoryRepository
	//ReactionRepository
	//NotificationRepository

}

type userRepository interface {
	CreateUser(user *entity.User) error
	UpdateUserPassword(user *entity.User) error
	//DeleteUser(userId int) error
	//GetUserByID(userId int) (entity.User, error)
	//CheckUserByEmail(email string) (bool, error)
	GetUserByEmail(email string) (*entity.User, error)
}

//type PostRepository interface {
//}
//
//type CommentRepository interface {
//}
//
//type CategoryRepository interface {
//}
//
//type PostCategoryRepository interface {
//}
//
//type ReactionRepository interface {
//}
//
//type NotificationRepository interface {
//}
