package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"SimpleForum/internal/config"
	"SimpleForum/internal/domain"
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
)

func (app *Application) CreatePost(userId int, title, content string, categories []string, myFile *entity.MyFile) error {
	var imageURL string
	titleValidation := checkContent(title)
	contentValidation := checkContent(content)

	if !(contentValidation && titleValidation) {
		return logger.ErrorWrapper("UseCase", "CreatePost", "The inserted title or content is not valid", domain.ErrNotValidContent)
	}

	if len(categories) == 0 {
		return logger.ErrorWrapper("UseCase", "CreatePost", "The Creating post cannot be created without categories", domain.ErrNoCategories)
	}

	if myFile != nil {
		// Check file size
		if myFile.FileHeader.Size > config.MaxImageSize {
			return logger.ErrorWrapper("UseCase", "CreatePost", "Image is too large (max 20MB)", domain.ErrLargeImageSize)
		}

		// Validate file type
		fileType := myFile.FileHeader.Header.Get("Content-Type")

		if !config.AllowedImageTypes[fileType] {
			return logger.ErrorWrapper("UseCase", "CreatePost", "Invalid file type. Only JPEG, PNG, and GIF are allowed.", domain.ErrInvalidImageType)
		}

		// Create the uploads directory if it doesn't exist
		if _, err := os.Stat(config.UploadDir); os.IsNotExist(err) {
			err := os.MkdirAll(config.UploadDir, os.ModePerm)
			if err != nil {
				return logger.ErrorWrapper("UseCase", "CreatePost", "Failed to create upload directory", err)
			}
		}

		// Save the file
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(myFile.FileHeader.Filename))
		filePath := filepath.Join(config.UploadDir, fileName)
		destFile, err := os.Create(filePath)
		if err != nil {
			return logger.ErrorWrapper("UseCase", "CreatePost", "Failed to save the file", err)
		}
		defer destFile.Close()

		_, err = destFile.ReadFrom(myFile.FileContent)
		if err != nil {
			return logger.ErrorWrapper("UseCase", "CreatePost", "Failed to save the file", err)
		}

		imageURL = strings.TrimPrefix(filePath, "uploads/")

	}

	postID, err := app.ServiceDB.CreatePost(userId, title, content, imageURL)
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
