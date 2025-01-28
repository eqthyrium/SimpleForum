package sqllite

import (
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
	"database/sql"
	"errors"
)

// We need to work with the Reaction table to retrieve all posts associated with a specific user. These posts should then be JOINed with the Post table to fetch all related details about the posts, and finally, return the results
func (rp *Repository) GetReactedPostsByCertainUser(userId int, reaction string) ([]entity.Posts, error) {

	var posts []entity.Posts
	var action string
	if reaction == "like" {
		action = "L"
	} else if reaction == "dislike" {
		action = "D"
	} else if reaction == "comment" {
		action = "C"
	}

	statement := `
		SELECT
		p.PostId, p.UserId, p.Title, p.Content, p.Image,
			p.LikeCount, p.DislikeCount, p.CreatedAt
		FROM
		Posts p
		INNER JOIN
		Reactions r ON p.PostId = r.PostId
		WHERE
		r.UserId = ? AND r.Action = ? AND CommentId IS NULL
		ORDER BY P.CreatedAt DESC
	`

	rows, err := rp.DB.Query(statement, userId, action)
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetReactedPostsByCertainUser", "The problem  in the getting  posts by certain user's liked reaction", err)
	}

	for rows.Next() {
		post := entity.Posts{}
		var image sql.NullString
		err := rows.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &image, &post.LikeCount, &post.DislikeCount, &post.CreatedAt)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "GetReactedPostsByCertainUser", "Failed to scan post row", err)
		}
		if image.Valid {
			post.Image = image.String
		} else {
			post.Image = ""
		}
		posts = append(posts, post)
	}

	if err := rows.Close(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetReactedPostsByCertainUser", "Failed to close the row of db", err)
	}

	return posts, nil

}

func (rp *Repository) GetReactionsNotificationOfPosts(userId int) ([]entity.Notifications, error) {

	notifications := []entity.Notifications{}
	statement := `
	SELECT u.Nickname, sbqq.Action, sbqq.PostId
	FROM (	
		SELECT r.UserId, r.Action, r.PostId
		FROM (
			SELECT  PostId
			FROM Posts
			WHERE UserId = ?
		) AS sbq
		INNER JOIN Reactions AS r ON sbq.PostId = r.PostId
		WHERE r.CommentId IS NULL AND r.UserId != ?
		ORDER BY r.CreatedAt DESC
		LIMIT 20
	) AS sbqq 
	INNER JOIN Users AS u ON sbqq.UserId = u.UserId
`
	rows, err := rp.DB.Query(statement, userId, userId)
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetReactionsNotificationOfPosts", "The problem  in the getting  notification for particular user", err)
	}

	for rows.Next() {
		notification := entity.Notifications{}
		err := rows.Scan(&notification.UserNickname, &notification.Action, &notification.PostId)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "GetReactionsNotificationOfPosts", "Failed to scan post row", err)
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Close(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetReactedPostsByCertainUser", "Failed to close the row of db", err)
	}

	return notifications, nil

}
func (rp *Repository) RetrieveExistenceOfReactionLD(userId int, identifier int, postOrComment string) (*entity.Reactions, error) {
	var statement string
	if postOrComment == "post" {
		statement = `SELECT Action FROM Reactions WHERE UserId = ? AND PostId = ? AND (Action = 'L' OR Action = 'D')`
	} else if postOrComment == "comment" {
		statement = `SELECT Action FROM Reactions WHERE UserId = ? AND CommentId = ? AND (Action = 'L' OR Action = 'D')`
	}

	row := rp.DB.QueryRow(statement, userId, identifier)
	// Map the result to the Reactions struct
	reaction := &entity.Reactions{}
	err := row.Scan(&reaction.Action)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil for both reaction and error to indicate no match
		} else {
			return nil, logger.ErrorWrapper(
				"Repository",
				"RetrieveExistenceOfReaction",
				"The problem within the process of getting the existence of reaction in the db",
				err,
			)
		}
	}

	return reaction, nil

}

func (rp *Repository) InsertReaction(userId, identifier int, postOrcomment, reaction string) error {
	var statement string
	var dbreaction string
	if reaction == "like" {
		dbreaction = "L"
	} else if reaction == "dislike" {
		dbreaction = "D"
	} else if reaction == "comment" {
		dbreaction = "C"
	} else {
		return logger.ErrorWrapper("Repository", "InsertReaction", "Invalid reaction type: must be 'like', 'dislike', or 'comment'", nil)
	}

	if postOrcomment == "post" {
		statement = `INSERT INTO Reactions (UserId,PostId,Action) VALUES (?,?,?) `
	} else if postOrcomment == "comment" {
		statement = `INSERT INTO Reactions (UserId,CommentId,Action) VALUES (?,?,?) `
	} else {
		return logger.ErrorWrapper("Repository", "InsertReaction", "Invalid postOrComment value: must be 'post' or 'comment'", nil)
	}

	_, err := rp.DB.Exec(statement, userId, identifier, dbreaction)
	if err != nil {
		return logger.ErrorWrapper("Repository", "InsertReaction", "The problem within the process of insertion of the reaction in db", err)
	}
	return nil
}

func (rp *Repository) DeleteReaction(userId, identifier int, postOrcomment, reaction string) error {
	var statement string
	var dbreaction string
	if reaction == "like" {
		dbreaction = "L"
	} else if reaction == "dislike" {
		dbreaction = "D"
	} else if reaction == "comment" {
		dbreaction = "C"
	}
	if postOrcomment == "post" {
		statement = `DELETE FROM Reactions WHERE UserId = ? AND PostId = ? AND Action = ?`
	} else if postOrcomment == "comment" {
		statement = `DELETE FROM Reactions WHERE UserId = ? AND CommentId = ? AND Action = ?`
	}
	_, err := rp.DB.Exec(statement, userId, identifier, dbreaction)
	if err != nil {
		return logger.ErrorWrapper("Repository", "DeleteReaction", "The problem within the process of deleting of the reaction in db", err)
	}
	return nil
}
