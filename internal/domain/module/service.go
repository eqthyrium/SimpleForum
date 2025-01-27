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
	GetCertainPostPage(postId int) (*entity.Posts, []entity.Commentaries, error)
	GetCertainPostInfo(postId int) (*entity.Posts, error)
	DeleteCertainPost(userId, postId int, role string) error
	EditCertainPost(userId, postId int, content string) error

	GetMyDislikedPosts(userId int) ([]entity.Posts, error)
	GetMyCommentedPosts(userId int) ([]entity.Posts, error)
}

type commentaries interface {
	GetLatestCommentaries(postId int) ([]entity.Commentaries, error)
	GetCertainCommentaryInfo(commentId int) (*entity.Commentaries, error)
	CreateCommentary(userId, postId int, content string) error
	EditCertainCommentary(userId, commentId int, content string) error
	DeleteCertainCommentary(userId, commentId int, role string) error

	GetComments(userId int) ([]entity.Commentaries, error)
}

type categories interface {
	GetAllCategories() ([]entity.Categories, error)
	CreateCategory(categoryName string) error
	DeleteCategory(categoryId int) error
}

type reactions interface {
	ExecutionOfReactionLD(userId, identifier int, postOrcomment, reactionLD string) error
}

type notifications interface {
	GetNotifications(userId int) ([]entity.Notifications, error)
}
