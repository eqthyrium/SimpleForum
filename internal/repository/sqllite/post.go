package sqllite

import (
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
	"database/sql"
	"strconv"
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
		statement = `SELECT * FROM Posts WHERE PostId IN (
														SELECT PostId FROM PostCategories WHERE CategoryId IN (?) 
														GROUP BY PostId
														) 
                  ORDER BY CreatedAt DESC`
		categoriesInt := make([]int, len(categories))
		for i := 0; i < len(categories); i++ {
			categoriesInt[i], err = strconv.Atoi(categories[i])
			if err != nil {
				return nil, logger.ErrorWrapper("Repository", "GetLatestAllPosts", "There is a problem with converting types from string to int", err)
			}

		}
		rows, err = rp.DB.Query(statement, categoriesInt)
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
