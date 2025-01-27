package usecase

import (
	"SimpleForum/internal/domain"
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
	"strconv"
)

func (app *Application) CreatePost(userId int, title, content string, categories []string) error {

	titleValidation := checkContent(title)
	contentValidation := checkContent(content)
	if !(contentValidation && titleValidation) {
		return logger.ErrorWrapper("UseCase", "CreatePost", "The inserted title or content is not valid", domain.ErrNotValidContent)
	}

	if len(categories) == 0 {
		return logger.ErrorWrapper("UseCase", "CreatePost", "The Creating post cannot be created without categories", domain.ErrNoCategories)
	}

	postID, err := app.ServiceDB.CreatePost(userId, title, content)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "CreatePost", "Failed to create a post", err)
	}

	for i := 0; i < len(categories); i++ {
		number, err := strconv.Atoi(categories[i])
		if err != nil {
			return logger.ErrorWrapper("UseCase", "CreatePost", "Failed to convert category to number", err)
		}
		err = app.ServiceDB.SetPostCategoryRelation(postID, number)
		if err != nil {
			return logger.ErrorWrapper("UseCase", "CreatePost", "Failed to create a post category relation", err)
		}
	}
	return nil
}
func (app *Application) GetLatestPosts(requestedCategories []string) ([]entity.Posts, error) {

	posts, err := app.ServiceDB.GetLatestAllPosts(requestedCategories)
	if err != nil {
		return nil, logger.ErrorWrapper("UseCase", "GetLatestPosts", "There is problem with getting all the recent posts from the db", err)
	}

	return posts, nil
}

func (app *Application) GetMyCreatedPosts(userId int) ([]entity.Posts, error) {

	posts, err := app.ServiceDB.GetPostsByCertainUser(userId)
	if err != nil {
		return nil, logger.ErrorWrapper("UseCase", "GetMyCreatedPosts", "There is problem with getting all my created posts from the db", err)
	}

	return posts, nil

}

func (app *Application) GetMyLikedPosts(userId int) ([]entity.Posts, error) {

	posts, err := app.ServiceDB.GetReactedPostsByCertainUser(userId, "like")
	if err != nil {
		return nil, logger.ErrorWrapper("UseCase", "GetMyLikedPosts", "There is problem with getting all my liked posts from the db", err)
	}

	return posts, nil
}

func (app *Application) GetCertainPostPage(postId int) (*entity.Posts, []entity.Commentaries, error) {

	post, err := app.ServiceDB.GetCertainPostInfo(postId)
	if err != nil {
		return nil, nil, logger.ErrorWrapper("UseCase", "GetCertainPostPage", "There is problem with getting certain post's info from the db", err)

	}
	commentaries, err := app.ServiceDB.GetCertainPostsCommentaries(postId)
	if err != nil {
		return nil, nil, logger.ErrorWrapper("UseCase", "GetCertainPostPage", "There is problem with getting all relating commentaries of the certain post from the db", err)
	}
	return post, commentaries, nil
}

func (app *Application) GetCertainPostInfo(postId int) (*entity.Posts, error) {
	post, err := app.ServiceDB.GetCertainPostInfo(postId)
	if err != nil {
		return nil, logger.ErrorWrapper("UseCase", "GetCertainPostInfo", "There is problem with getting certain post", err)
	}
	return post, nil
}

func (app *Application) DeleteCertainPost(userId, postId int, role string) error {
	validation, err := app.ServiceDB.ValidateOfExistenceCertainPost(userId, postId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "DeleteCertainPost", "There is problem with validating certain post", err)
	}
	if !validation && role == "User" {
		return logger.ErrorWrapper("UseCase", "DeleteCertainPost", "There is no such post of that user", domain.ErrPostNotFound)
	}

	err = app.ServiceDB.DeleteCertainPost(postId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "DeleteCertainPost", "There is problem with deleting certain post", err)
	}

	return nil

}

func (app *Application) EditCertainPost(userId, postId int, content string) error {

	validation, err := app.ServiceDB.ValidateOfExistenceCertainPost(userId, postId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "EditCertainPost", "There is problem with validating certain post", err)
	}
	if !validation {
		return logger.ErrorWrapper("UseCase", "EditCertainPost", "There is no such post of that user", domain.ErrPostNotFound)
	}

	validation = checkContent(content)
	if !validation {
		return logger.ErrorWrapper("UseCase", "EditCertainPost", "There is no content", domain.ErrNotValidContent)
	}

	err = app.ServiceDB.UpdateEditedPost(userId, postId, content)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "EditedCertainPost", "Failed to update the content of the post", err)
	}

	return nil
}

// ------------------------- Asem code is started here below
func (app *Application) GetMyDislikedPosts(userId int) ([]entity.Posts, error) {
	posts, err := app.ServiceDB.GetReactedPostsByCertainUser(userId, "dislike")
	if err != nil {
		return nil, logger.ErrorWrapper("UseCase", "GetMyLikedPosts", "There is problem with getting all my liked posts from the db", err)
	}

	return posts, nil
}

func (app *Application) GetMyCommentedPosts(userId int) ([]entity.Posts, error) {
	posts, err := app.ServiceDB.GetMyCommentedPosts(userId)
	if err != nil {
		return nil, logger.ErrorWrapper("UseCase", "GetLatestPosts", "There is problem with getting all the recent posts from the db", err)
	}
	return posts, nil
}
