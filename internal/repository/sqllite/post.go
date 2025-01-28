package sqllite

import (
	"SimpleForum/internal/domain"
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func (rp *Repository) CreatePost(userId int, title, content, imageURL string) (int, error) {

	var postId int64
	if imageURL == "" {
		statement := `INSERT INTO Posts (UserId, Title, Content) VALUES(?,?,?)`
		index, err := rp.DB.Exec(statement, userId, title, content)

		if err != nil {
			return -1, logger.ErrorWrapper("Repository", "CreatePost", "The problem within the process of creation of the user in db", err)
		}
		postId, err = index.LastInsertId()
		if err != nil {
			return -1, logger.ErrorWrapper("Repository", "CreatePost", "The problem within the process of getting the last ID of Posts table  in db", err)
		}
	} else {
		statement := `INSERT INTO Posts (UserId, Title, Content,Image) VALUES(?,?,?,?)`
		index, err := rp.DB.Exec(statement, userId, title, content, imageURL)

		if err != nil {
			return -1, logger.ErrorWrapper("Repository", "CreatePost", "The problem within the process of creation of the user in db", err)
		}
		postId, err = index.LastInsertId()
		if err != nil {
			return -1, logger.ErrorWrapper("Repository", "CreatePost", "The problem within the process of getting the last ID of Posts table  in db", err)
		}
	}

	return int(postId), nil
}

func (rp *Repository) GetLatestAllPosts(categories []string) ([]entity.Posts, error) {
	var statement string
	var err error
	var rows *sql.Rows
	var posts []entity.Posts

	if len(categories) == 0 {
		statement = `SELECT * FROM Posts ORDER BY CreatedAt DESC`
		rows, err = rp.DB.Query(statement)
	} else {
		categoriesString := strings.Join(categories, ",")
		statement = fmt.Sprintf("SELECT * FROM Posts WHERE PostId IN ( SELECT PostId FROM PostCategories WHERE CategoryId IN (%s) GROUP BY PostId) ORDER BY CreatedAt DESC", categoriesString)
		rows, err = rp.DB.Query(statement)
	}

	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetLatestAllPosts", "Failed to execute query for posts", err)
	}

	for rows.Next() {
		post := entity.Posts{}
		var image sql.NullString

		err := rows.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &image, &post.LikeCount, &post.DislikeCount, &post.CreatedAt)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "GetLatestAllPosts", "Failed to scan post row", err)
		}

		if image.Valid {
			post.Image = image.String
		} else {
			post.Image = ""
		}

		posts = append(posts, post)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetLatestAllPosts", "Error occurred during rows iteration", err)
	}

	if err := rows.Close(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetLatestAllPosts", "Failed to close the row of db", err)
	}
	return posts, nil
}

// Here we have to return all related posts to the certain User
func (rp *Repository) GetPostsByCertainUser(userId int) ([]entity.Posts, error) {
	var posts []entity.Posts
	statement := `
	SELECT *  FROM Posts 	
	          WHERE UserId = ? 
			ORDER BY CreatedAt DESC
	`
	rows, err := rp.DB.Query(statement, userId)
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetPostsByCertainUser", "The problem  in the getting  posts by certain user", err)
	}

	for rows.Next() {
		post := entity.Posts{}
		var image sql.NullString
		err := rows.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &image, &post.LikeCount, &post.DislikeCount, &post.CreatedAt)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "GetPostsByCertainUser", "Failed to scan post row", err)
		}
		if image.Valid {
			post.Image = image.String
		} else {
			post.Image = ""
		}
		posts = append(posts, post)
	}

	if err := rows.Close(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetPostsByCertainUser", "Failed to close the row of db", err)
	}
	return posts, nil
}

func (rp *Repository) GetCertainPostInfo(postId int) (*entity.Posts, error) {

	var statement string = `Select * FROM Posts WHERE PostId = ?`

	row := rp.DB.QueryRow(statement, postId)

	post := &entity.Posts{}
	var image sql.NullString

	err := row.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &image, &post.LikeCount, &post.DislikeCount, &post.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, logger.ErrorWrapper("Repository", "GetCertainPostInfo", "There is no such  post in the db", domain.ErrPostNotFound)
		} else {
			return nil, logger.ErrorWrapper("Repository", "GetCertainPostInfo", "The problem within the process of getting of the particular post by its postId in db", err)
		}
	}

	if image.Valid {
		post.Image = image.String
	} else {
		post.Image = ""
	}

	return post, nil
}

func (rp *Repository) UpdateReactionOfPost(postId int, reaction, operation string) error {

	var statement string

	if reaction == "like" {
		if operation == "increment" {
			statement = `UPDATE Posts SET LikeCount = LikeCount + 1 WHERE PostId = ?`
		} else if operation == "decrement" {
			statement = `UPDATE Posts SET LikeCount = LikeCount - 1 WHERE PostId = ? and LikeCount > 0`
		}
	} else if reaction == "dislike" {
		if operation == "increment" {
			statement = `UPDATE Posts SET DislikeCount = DislikeCount + 1 WHERE PostId = ?`
		} else if operation == "decrement" {
			statement = `UPDATE Posts SET DislikeCount = DislikeCount - 1 WHERE PostId = ? and DislikeCount > 0`
		}
	}

	_, err := rp.DB.Exec(statement, postId)
	if err != nil {
		return logger.ErrorWrapper(
			"Repository",
			"UpdateReactionOfPost",
			"Failed to increment/decrement reaction counter for the post in the database",
			err,
		)
	}

	return nil
}

func (rp *Repository) UpdateEditedPost(userId, postId int, content string) error {
	statement := `
		UPDATE Posts
		SET Content = ?
		WHERE PostId = ? AND UserId = ?
	`

	result, err := rp.DB.Exec(statement, content, postId, userId)
	if err != nil {
		return logger.ErrorWrapper("Repository", "UpdateEditedPost", "Failed to update post:", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return logger.ErrorWrapper("Repository", "UpdateEditedPost", "Failed to retrieve affected rows:", err)

	}
	if rowsAffected == 0 {
		return logger.ErrorWrapper("Repository", "UpdateEditedPost", "No rows were updated; post_id or user_id might not exist", err)
	}

	return nil
}

func (rp *Repository) DeleteCertainPost(postId int) error {

	statement := `
		PRAGMA foreign_keys = ON;	
        DELETE FROM Posts WHERE PostId = ?;
`

	_, err := rp.DB.Exec(statement, postId)
	if err != nil {
		return logger.ErrorWrapper("Repository", "DeleteCertainPost", "The problem within the process of deleting of the post in db", err)
	}
	return nil
}

func (rp *Repository) ValidateOfExistenceCertainPost(userId, postId int) (bool, error) {
	var exists bool

	var statement string = `
		SELECT EXISTS(
			SELECT 1 FROM Posts WHERE PostId = ? AND UserId = ?
		)
	`

	err := rp.DB.QueryRow(statement, postId, userId).Scan(&exists)

	if err != nil {
		return false, logger.ErrorWrapper("Repository", "ValidateOfExistenceCertainPost", "Error checking existence of the post", err)
	}

	return exists, nil
}

//--------------------

func (rp *Repository) GetMyCommentedPosts(userId int) ([]entity.Posts, error) {
	var posts []entity.Posts

	// Запрос с JOIN для получения информации о постах, на которые пользователь оставил комментарии
	stmt := `SELECT DISTINCT Posts.PostId, Posts.UserId, Posts.Title, Posts.Content, Posts.Image, Posts.LikeCount, Posts.DislikeCount, Posts.CreatedAt
	FROM Commentaries
	JOIN Posts ON Posts.PostId = Commentaries.PostId
	WHERE Commentaries.UserId = ?`

	// Выполнение запроса
	rows, err := rp.DB.Query(stmt, userId)
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetMyCommentedPosts", "Problem getting posts commented by user", err)
	}
	defer rows.Close() // безопасное закрытие ресурсов

	// Обработка результата запроса
	for rows.Next() {
		post := entity.Posts{}
		var image sql.NullString
		err := rows.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &image, &post.LikeCount, &post.DislikeCount, &post.CreatedAt)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "GetMyCommentedPosts", "Failed to scan post row", err)
		}
		if image.Valid {
			post.Image = image.String
		} else {
			post.Image = ""
		}
		posts = append(posts, post)
	}

	// Проверка на ошибки после завершения обработки строк
	if err := rows.Err(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetMyCommentedPosts", "Error iterating over rows", err)
	}

	return posts, nil
}
