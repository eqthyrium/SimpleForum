package sqllite

import (
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
)

func (rp *Repository) CreateCommentary(userId, postId int, content string) error {

	statement := `INSERT INTO Commentaries (UserId, PostId, Content) VALUES (?, ?,?)`
	_, err := rp.DB.Exec(statement, userId, postId, content)
	if err != nil {
		return logger.ErrorWrapper("Repository", "CreateCommentary", "Failed to execute statement", err)
	}
	return nil
}

func (rp *Repository) GetCertainPostsCommentaries(postId int) ([]entity.Commentaries, error) {

	statement := `SELECT * FROM Commentaries WHERE PostId = ? ORDER BY CreatedAt DESC`

	rows, err := rp.DB.Query(statement, postId)

	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetCertainPostsCommentaries", "Failed to execute query for commentaries", err)
	}
	var commentaries []entity.Commentaries

	for rows.Next() {
		commentary := entity.Commentaries{}
		err := rows.Scan(&commentary.CommentId, &commentary.PostId, &commentary.UserId, &commentary.Content, &commentary.LikeCount, &commentary.DislikeCount, &commentary.CreateAt)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "GetCertainPostsCommentaries", "Failed to scan commentaries row", err)
		}
		commentaries = append(commentaries, commentary)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetCertainPostsCommentaries", "Error occurred during rows iteration", err)
	}

	if err := rows.Close(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetCertainPostsCommentaries", "Failed to close the row of db", err)
	}
	return commentaries, nil
}

func (rp *Repository) UpdateReactionOfCommentary(commentId int, reaction, operation string) error {

	var statement string

	if reaction == "like" {
		if operation == "increment" {
			statement = `UPDATE Commentaries SET LikeCount = LikeCount + 1 WHERE CommentId = ?`
		} else if operation == "decrement" {
			statement = `UPDATE Commentaries SET LikeCount = LikeCount - 1 WHERE CommentId = ? and LikeCount > 0`
		}
	} else if reaction == "dislike" {
		if operation == "increment" {
			statement = `UPDATE Commentaries SET DislikeCount = DislikeCount + 1 WHERE CommentId = ?`
		} else if operation == "decrement" {
			statement = `UPDATE Commentaries SET DislikeCount = DislikeCount - 1 WHERE CommentId = ? and DislikeCount > 0`
		}
	}

	_, err := rp.DB.Exec(statement, commentId)
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
