package module

import "SimpleForum/internal/domain/entity"

type HttpModule interface {
	authentication
	posts
	commentaries
	categories
	reactions
	notifications
}

type authentication interface {
	SignUp(nickname, email, password, oauth string) error
	LogIn(email, password, oauth string) (string, error)
}

type posts interface {
	CreatePost(userId int, title, content string, categories []string) error
	GetLatestPosts(requestedCategories []string) ([]entity.Posts, error)
	GetMyCreatedPosts(userId int) ([]entity.Posts, error)
	GetMyLikedPosts(userId int) ([]entity.Posts, error)
	GetMyCommentedPosts(userId int) ([]entity.Posts, error)
	GetCertainPostPage(postId int) (*entity.Posts, []entity.Commentaries, error)
}

type commentaries interface {
	GetLatestCommentaries(postId int) ([]entity.Commentaries, error)
	CreateCommentary(userId, postId int, content string) error
	GetCertainPostsCommentaries(postId int) ([]entity.Commentaries, error)
}

type categories interface {
	GetAllCategories() ([]entity.Categories, error)
}

type reactions interface {
	ExecutionOfReactionLD(userId, identifier int, postOrcomment, reactionLD string) error
}

type notifications interface {
	GetNotifications(userId int) ([]entity.Notifications, error)
}
