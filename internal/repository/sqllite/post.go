package sqllite

import (
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
	"database/sql"
	"fmt"
	"strings"
)

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
		fmt.Println("Layer:Repository", "the list of selected categories from the client side:", categoriesString)
		statement = fmt.Sprintf("SELECT * FROM Posts WHERE PostId IN ( SELECT PostId FROM PostCategories WHERE CategoryId IN (%s) GROUP BY PostId) ORDER BY CreatedAt DESC", categoriesString)
		rows, err = rp.DB.Query(statement)
	}

	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetLatestAllPosts", "Failed to execute query for posts", err)
	}
	defer rows.Close()

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

	return posts, nil
}

// Here we have to return all related posts to the certain User
//func (rp *Repository) GetPostsByCertainUser(userId int) ([]entity.Posts, error) {
//
//}
