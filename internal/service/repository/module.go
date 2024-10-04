package repository

import (
	"SimpleForum/internal/domain"
)

type dbModule interface {
	UserRepository
	//PostRepository
	//CommentRepository
	//CategoryRepository
	//PostCategoryRepository
	//ReactionRepository
	//NotificationRepository

}

type ServiceRepository struct {
	Repo dbModule
}

func NewServiceRepository(repoObject dbModule) *ServiceRepository {
	return &ServiceRepository{Repo: repoObject}
}

type UserRepository interface {
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
