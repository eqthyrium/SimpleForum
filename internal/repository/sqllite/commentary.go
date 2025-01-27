package sqllite

import (
	"SimpleForum/internal/domain"
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
	"database/sql"
	"errors"
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

func (rp *Repository) GetCertainCommentaryInfo(commentId int) (*entity.Commentaries, error) {

	var statement string = `Select * FROM Commentaries WHERE CommentId = ?`

	row := rp.DB.QueryRow(statement, commentId)

	commentary := &entity.Commentaries{}

	err := row.Scan(&commentary.CommentId, &commentary.PostId, &commentary.UserId, &commentary.Content, &commentary.LikeCount, &commentary.DislikeCount, &commentary.CreateAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, logger.ErrorWrapper("Repository", "GetCertainCommentaryInfo", "There is no such  post in the db", domain.ErrCommentaryNotFound)
		} else {
			return nil, logger.ErrorWrapper("Repository", "GetCertainCommentaryInfo", "The problem within the process of getting of the particular post by its postId in db", err)
		}
	}

	return commentary, nil
}

func (rp *Repository) UpdateEditedCommentary(userId, commentId int, content string) error {
	statement := `
		UPDATE Commentaries
		SET Content = ?
		WHERE CommentId = ? AND UserId = ?
	`

	result, err := rp.DB.Exec(statement, content, commentId, userId)
	if err != nil {
		return logger.ErrorWrapper("Repository", "UpdateEditedCommentary", "Failed to update post:", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return logger.ErrorWrapper("Repository", "UpdateEditedCommentary", "Failed to retrieve affected rows:", err)

	}
	if rowsAffected == 0 {
		return logger.ErrorWrapper("Repository", "UpdateEditedCommentary", "No rows were updated; post_id or user_id might not exist", err)
	}

	return nil
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

func (rp *Repository) DeleteCertainCommentary(commentId int) error {

	statement := `
		PRAGMA foreign_keys = ON;
		DELETE FROM Commentaries WHERE CommentId = ?
	`

	_, err := rp.DB.Exec(statement, commentId)
	if err != nil {
		return logger.ErrorWrapper("Repository", "DeleteCertainCommentary", "The problem within the process of deleting of the commentary in db", err)
	}
	return nil
}

func (rp *Repository) ValidateOfExistenceCertainCommentary(userId, commentId int) (bool, error) {
	var exists bool

	var statement string = `
		SELECT EXISTS(
			SELECT 1 FROM Commentaries WHERE  CommentId= ? AND UserId = ?
		)
	`

	err := rp.DB.QueryRow(statement, commentId, userId).Scan(&exists)

	if err != nil {
		return false, logger.ErrorWrapper("Repository", "ValidateOfExistenceCertainPost", "Error checking existence of the post", err)
	}

	return exists, nil
}

//---------------------------

func (rp *Repository) GetComments(UserId int) ([]entity.Commentaries, error) {
	stmt := `SELECT * FROM Commentaries WHERE UserId = ?`
	rows, err := rp.DB.Query(stmt, UserId)
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetComments", "Failed to execute query for commentaries", err)
	}
	var commentaries []entity.Commentaries

	for rows.Next() {
		commentary := entity.Commentaries{}
		err := rows.Scan(&commentary.CommentId, &commentary.PostId, &commentary.UserId, &commentary.Content, &commentary.LikeCount, &commentary.DislikeCount, &commentary.CreateAt)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "GetCertainPostsGetCommentsCommentaries", "Failed to scan commentaries row", err)
		}
		commentaries = append(commentaries, commentary)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetComments", "Error occurred during rows iteration", err)
	}

	if err := rows.Close(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetComments", "Failed to close the row of db", err)
	}

	return commentaries, nil
}
