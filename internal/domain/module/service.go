package module

import "SimpleForum/internal/domain/entity"

type HttpModule interface {
	SignUp(nickname, email, password, oauth string) error
	LogIn(email, password, oauth string) (string, error)
	GetAllCategories() ([]entity.Categories, error)
	GetLatestPosts(requestedCategories []string) ([]entity.Posts, error)
	GetLatestCommentaries(postId int) ([]entity.Commentaries, error)
	//Login(email, password string) (err error)
	// Here we write what kind of services can be used in the http handler
}
