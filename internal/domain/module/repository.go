package module

import (
	"SimpleForum/internal/domain/entity"
)

type DbModule interface {
	userRepository
	postRepository
	//commentRepository
	categoryRepository
	//PostCategoryRepository
	//ReactionRepository
	//NotificationRepository

}

type userRepository interface {
	CreateUser(user *entity.Users) error
	UpdateUserPassword(user *entity.Users) error
	//DeleteUser(userId int) error
	//GetUserByID(userId int) (entity.User, error)
	//CheckUserByEmail(email string) (bool, error)
	GetUserByEmail(email string) (*entity.Users, error)
}

type postRepository interface {
	GetLatestAllPosts(categories []string) ([]entity.Posts, error)
	GetPostsByCertainUser(userId int) ([]entity.Posts, error)
}

//type commentRepository interface {
//}

type categoryRepository interface {
	GetAllCategories() ([]entity.Categories, error)
}

//
//type PostCategoryRepository interface {
//}

type ReactionRepository interface {
	GetReactedPostsByCertainUser(userId int, reaction string) ([]entity.Posts, error)
}

//
//type NotificationRepository interface {
//}
