package module

import (
	"SimpleForum/internal/domain/entity"
)

type DbModule interface {
	userRepository
	postRepository
	commentRepository
	categoryRepository
	reactionRepository
	postCategoryRepository
	// PostCategoryRepository
	// ReactionRepository
	// NotificationRepository
}

type userRepository interface {
	CreateUser(user *entity.Users) error
	UpdateUserPassword(user *entity.Users) error
	// DeleteUser(userId int) error
	// GetUserByID(userId int) (entity.User, error)
	// CheckUserByEmail(email string) (bool, error)
	GetUserByEmail(email string) (*entity.Users, error)
}

type postRepository interface {
	CreatePost(userId int, title, content string) (int, error)
	GetLatestAllPosts(categories []string) ([]entity.Posts, error)
	GetMyCommentedPosts(userId int) ([]entity.Posts, error)
	GetPostsByCertainUser(userId int) ([]entity.Posts, error)
	GetCertainPostInfo(postId int) (*entity.Posts, error)
	UpdateReactionOfPost(postId int, reaction, operation string) error
}

type commentRepository interface {
	GetComments(UserId int) ([]entity.Commentaries, error)
	GetCertainPostsCommentaries(postId int) ([]entity.Commentaries, error)
	CreateCommentary(userId, postId int, content string) error
	UpdateReactionOfCommentary(commentId int, reaction, operation string) error
}

type categoryRepository interface {
	AddCategory(categoryName string)([]entity.Categories, error)
	DeleteCategory(categoryId int)([]entity.Categories, error)
	GetAllCategories() ([]entity.Categories, error)
}

type postCategoryRepository interface {
	SetPostCategoryRelation(postId, categoryId int) error
}

type reactionRepository interface {
	GetReactedPostsByCertainUser(userId int, reaction string) ([]entity.Posts, error) // Azamat
	GetReactionsNotificationOfPosts(userId int) ([]entity.Notifications, error)
	RetrieveExistenceOfReactionLD(userId int, identifier int, postOrComment string) (*entity.Reactions, error)
	InsertReaction(userId, identifier int, postOrcomment, reaction string) error
	DeleteReaction(userId, identifier int, postOrComment, reaction string) error
}

//
//type NotificationRepository interface {
//}
