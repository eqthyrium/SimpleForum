package repository

import (
	"SimpleForum/internal/domain"
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
	CreateUser(user *domain.User) error
	//UpdateUser()error
	//DeleteUser(userId int) error
	//GetUserByID(userId int) (domain.User, error)
	//CheckUserByEmail(email string) (bool, error)
	GetUserByEmail(email string) (*domain.User, error)
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
