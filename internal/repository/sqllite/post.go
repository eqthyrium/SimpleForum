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

func (rp *Repository) CreatePost(post *entity.Posts) error {
	statement := `
	INSERT INTO Posts
	UserId,
	Title,
	Content,
	Image
	VALUES (?, ?, ?, ?)
	`
	_, err := rp.DB.Exec(statement, post.UserId, post.Title, post.Content, post.Image)
	if err != nil {
		return logger.ErrorWrapper("Repository", "CreatePost", "The problem within the process of creation of the post in db", err)
	}
	return nil
}

func (rp *Repository) DeletePost(PostId int) error {
	statement := `
	DELETE  
	FROM Posts
	WHERE PostId = ?
	`
	_, err := rp.DB.Exec(statement, PostId)
	if err != nil {
		return logger.ErrorWrapper("Repository", "DeletePost", "The problem within the process of delete the post in db", err)
	}
	return nil
}

func (rp *Repository) LikePost(UserID, PostID int) error {
	if rp.IsPostLiked(UserID, PostID) {
		query := `
		UPDATE Posts
			SET LikeCount = LikeCount - 1 
			WHERE PostId = ? and UserId = ?
		`
		_, err := rp.DB.Exec(query, PostID, UserID)
		if err != nil {
			return logger.ErrorWrapper("Repository", "LikePost", "The problem in the LikePost function (query in Posts table)", err)
		}
		query = `
		DELETE FROM Reactions
		WHERE UserId = ? and PosrId = ? and Action = 'L'
		`
		_, err = rp.DB.Exec(query, UserID, PostID)
		if err != nil {
			return logger.ErrorWrapper("Repository", "LikePost", "The problem in the LikePost function (query in Reaction table)", err)
		}
		return nil
	} else if rp.IsPostDisliked(UserID, PostID) {
		query := `
		UPDATE Posts
			SET DislikeCount = DislikeCount - 1
			WHERE PostId = ? and UserId = ?
		`
		_, err := rp.DB.Exec(query, PostID, UserID)
		if err != nil {
			return logger.ErrorWrapper("Repository", "LikePost", "The problem in the LikePost function (query in Posts table)", err)
		}
		query = `
		DELETE FROM Reactions
		WHERE UserId = ? and PostId = ? and Action = 'D'
		`
		_, err = rp.DB.Exec(query, UserID, PostID)
		if err != nil {
			return logger.ErrorWrapper("Reposotpry", "LikePost", "The problem in the LikePost function (query in Reaction table)", err)
		}
		return nil
	} else if !rp.IsPostDisliked(UserID, PostID) && !rp.IsPostLiked(UserID, PostID) {
		query := `
		INSERT INTO Posts 
			PostId = ?,
			UserId = ?,
			LikeCount = LikeCount + 1 
			VALUES (?, ?) 
		`
		_, err := rp.DB.Exec(query, PostID, UserID)
		if err != nil {
			return logger.ErrorWrapper("Repository", "LikePost", "The problem is to insert like into posts tables", err)
		}
		query = `
		INSERT INTO Reactions
			PostId = ?,
			UserId = ?
			Action = 'L'
		`
		_, err = rp.DB.Exec(query, PostID, UserID)
		if err != nil {
			return logger.ErrorWrapper("Repository", "LikePost", "The problem is to insert into reation tables", err)
		}
	}
	return nil
}

func (rp *Repository) DislikePost(UserID, PostID int) error {
	if rp.IsPostLiked(UserID, PostID) {
		query := `
		UPDATE Posts
		LikeCount = LikeCount - 1
		WHERE UserId = ? and PostId = ?
		`
		_, err := rp.DB.Exec(query, UserID, PostID)
		if err != nil {
			return logger.ErrorWrapper("Repository", "DisLikePost", "The problem is to update post", err)
		}
		query = `
		DELETE FROM Reactions
		WHERE UserId = ? and PostId = ? and Action = 'L'
		`
		_, err = rp.DB.Exec(query, UserID, PostID)
		if err != nil {
			return logger.ErrorWrapper("Repository", "DisLikePost", "The problem is delete post reactions", err)
		}
		return nil
	} else if rp.IsPostDisliked(UserID, PostID) {
		return nil
	} else if !rp.IsPostLiked(UserID, PostID) && !rp.IsPostDisliked(UserID, PostID) {
		query := `
		INSERT INTO Posts
		UserId = ?
		PostId = ?
		DislikeCount = DislikeCount + 1
		`
		_, err := rp.DB.Exec(query, UserID, PostID)
		if err != nil {
			return logger.ErrorWrapper("Repository", "DisLikePost", "The problem is insert into Posts tables", err)
		}
		query = `
		INSERT INTO Reactions
		UserId = ?,
		PostId = ?
		Action = 'D'
		VALUES = (?, ?)
		`
		_, err = rp.DB.Exec(query, UserID, PostID)
		if err != nil {
			return logger.ErrorWrapper("Repository", "DisLikePost", "The problem is insert into Reactions tables", err)
		}
	}
	return nil
}

func (rp *Repository) IsPostLiked(UserID int, PostID int) bool {
	var answer bool
	query := `
	EXIST (SELECT *
			FROM Reactions
			WHERE UserId = ? and PostId = ? and Actions = 'L')
	`
	err := rp.DB.QueryRow(query, UserID, PostID).Scan(&answer)
	if err != nil {
		logger.ErrorWrapper("Repository", "Isliked", "The problem in the IsLiked function", err)
		return answer
	}
	return answer
}

func (rp *Repository) IsPostDisliked(UserID int, PostID int) bool {
	var answer bool
	query := `
	EXIST (SELECT *
			FROM Reactions
			WHERE UserId = ? and PostId = ? and Actions = 'D')
	`
	err := rp.DB.QueryRow(query, UserID, PostID).Scan(&answer)
	if err != nil {
		logger.ErrorWrapper("Repository", "Isliked", "The problem in the IsLiked function", err)
		return answer
	}
	return answer
}
func (rp *Repository) GetPostsByCertainUser(UserId int) ([]entity.Posts, error) {
	var posts []entity.Posts
	query := `
	SELECT * 
	FROM Posts
	WHERE UserId = ? 
	ORDER BY CreatedAt DESC
	`
	rows, err := rp.DB.Query(query, UserId)
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetPostsByCertainUser", "The problem is get posts by user", err)
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
	return posts, nil
}

func (rp *Repository) GetReactedPostsByCertainUser(UserId int) ([]entity.Posts, error) {
	var posts []entity.Posts
	query := `
    WITH Tb1 AS (
        SELECT PostId
        FROM Reactions
        WHERE UserId = ? AND Action IN ('L')
    )
    SELECT * 
    FROM Posts
    WHERE PostId IN (SELECT PostId FROM Tb1)
    ORDER BY CreatedAt DESC
    `
	rows, err := rp.DB.Query(query, UserId)
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetReactedPostsByCertainUser", "The problem is get posts by user", err)
	}
	defer rows.Close()

	for rows.Next() {
		post := entity.Posts{}
		err := rows.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &post.Image, &post.LikeCount, &post.DislikeCount, &post.CreatedAt)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "GetReactedPostsByCertainUser", "Failed to scan post row", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetReactedPostsByCertainUser", "Error occurred during rows iteration", err)
	}

	return posts, nil
}
