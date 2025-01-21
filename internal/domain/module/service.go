package module

import "SimpleForum/internal/domain/entity"

type HttpModule interface {
	Authentication
	Posts
	Commentaries
	Categories
}

type Authentication interface {
	SignUp(nickname, email, password, oauth string) error
	LogIn(email, password, oauth string) (string, error)
}

type Posts interface {
	GetLatestPosts(requestedCategories []string) ([]entity.Posts, error)
}

type Commentaries interface {
	GetLatestCommentaries(postId int) ([]entity.Commentaries, error)
}

type Categories interface {
	GetAllCategories() ([]entity.Categories, error)
}
