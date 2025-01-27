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
	CreatePost(userId int, title, content string) (int, error)
	GetLatestAllPosts(categories []string) ([]entity.Posts, error)
	GetPostsByCertainUser(userId int) ([]entity.Posts, error)
	GetCertainPostInfo(postId int) (*entity.Posts, error)
	UpdateReactionOfPost(postId int, reaction, operation string) error
	UpdateEditedPost(userId, postId int, content string) error
	DeleteCertainPost(postId int) error
	ValidateOfExistenceCertainPost(userId, postId int) (bool, error)

	GetMyCommentedPosts(userId int) ([]entity.Posts, error)
}

type commentRepository interface {
	CreateCommentary(userId, postId int, content string) error
	GetCertainPostsCommentaries(postId int) ([]entity.Commentaries, error)
	GetCertainCommentaryInfo(commentId int) (*entity.Commentaries, error)
	UpdateReactionOfCommentary(commentId int, reaction, operation string) error
	UpdateEditedCommentary(userId, commentId int, content string) error
	DeleteCertainCommentary(commentId int) error
	ValidateOfExistenceCertainCommentary(userId, commentId int) (bool, error)

	GetComments(UserId int) ([]entity.Commentaries, error)
}

type categoryRepository interface {
	GetAllCategories() ([]entity.Categories, error)
	CreateCategory(categoryName string) error
	DeleteCategory(categoryId int) error
}

type postCategoryRepository interface {
	SetPostCategoryRelation(postId, categoryId int) error
	GetPostIDsOfCertainCategory(categoryId int) ([]int, error)
	GetCategoriesOfCertainPost(postId int) ([]int, error)
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
