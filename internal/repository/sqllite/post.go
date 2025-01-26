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

func (rp *Repository) CreatePost(userId int, title, content string) (int, error) {

	statement := `INSERT INTO Posts (UserId, Title, Content) VALUES(?,?,?)`
	index, err := rp.DB.Exec(statement, userId, title, content)
	if err != nil {
		return -1, logger.ErrorWrapper("Repository", "CreatePost", "The problem within the process of creation of the user in db", err)
	}
	postId, err := index.LastInsertId()
	if err != nil {
		return -1, logger.ErrorWrapper("Repository", "CreatePost", "The problem within the process of getting the last ID of Posts table  in db", err)
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
		err := rows.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &post.Image, &post.LikeCount, &post.DislikeCount, &post.CreatedAt)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "GetLatestAllPosts", "Failed to scan post row", err)
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
		err := rows.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &post.Image, &post.LikeCount, &post.DislikeCount, &post.CreatedAt)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "GetPostsByCertainUser", "Failed to scan post row", err)
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

	err := row.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &post.Image, &post.LikeCount, &post.DislikeCount, &post.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, logger.ErrorWrapper("Repository", "GetCertainPostInfo", "There is no such  post in the db", domain.ErrPostNotFound)
		} else {
			return nil, logger.ErrorWrapper("Repository", "GetCertainPostInfo", "The problem within the process of getting of the particular post by its postId in db", err)
		}
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
