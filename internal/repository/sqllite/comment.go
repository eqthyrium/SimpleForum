package sqllite

import (
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
)

func (rp *Repository) CreateComment(comment *entity.Commentaries) error {
	query := `
	INSERT INTO Commentaries
	(PostId,
	UserId,
	Content)
	Values = (?, ?, ?)  
	`
	_, err := rp.DB.Exec(query, comment.PostId, comment.UserId, comment.Content)
	if err != nil {
		return logger.ErrorWrapper("Repository", "CreateComment", "The Problem is to create comment in db", err)
	}
	return nil
}

func (rp *Repository) UpdateComment(comment *entity.Commentaries) error {
	query := `
	UPDATE Commentaries
	SET Content = ?
	WHERE CommentId = ?
	`
	_, err := rp.DB.Exec(query, comment.Content, comment.CommentId)
	if err != nil {
		return logger.ErrorWrapper("Repository", "UpdateComment", "The problem within the process update comment in DB", err)
	}
	return nil
}

func (rp *Repository) DeleteComment(CommentId int) error {
	query := `
	DELETE FROM Commentaries
	WHERE CommentId = ?
	`
	_, err := rp.DB.Exec(query, CommentId)
	if err != nil {
		return logger.ErrorWrapper("Repository", "DeleteComment", "The problem within the process delete comment in DB", err)
	}
	return nil
}

func (rp *Repository) LikeComment(UserId, CommentId int) error {
	if rp.IsCommentLiked(UserId, CommentId) {
		query := `
        UPDATE Commentaries
        SET LikeCount = LikeCount - 1 
        WHERE CommentId = ? AND UserId = ?
        `
		_, err := rp.DB.Exec(query, CommentId, UserId)
		if err != nil {
			return logger.ErrorWrapper("Repository", "LikeComment", "The problem in the LikeComment function (query in Commentaries table)", err)
		}
		query = `
        DELETE FROM Reactions
        WHERE UserId = ? AND CommentId = ? AND Action = 'L'
        `
		_, err = rp.DB.Exec(query, UserId, CommentId)
		if err != nil {
			return logger.ErrorWrapper("Repository", "LikeComment", "The problem in the LikeComment function (query in Reactions table)", err)
		}
		return nil
	} else if rp.IsCommentDisliked(UserId, CommentId) {
		query := `
        UPDATE Commentaries
        SET DislikeCount = DislikeCount - 1
        WHERE CommentId = ? AND UserId = ?
        `
		_, err := rp.DB.Exec(query, CommentId, UserId)
		if err != nil {
			return logger.ErrorWrapper("Repository", "LikeComment", "The problem in the LikeComment function (query in Commentaries table)", err)
		}
		query = `
        DELETE FROM Reactions
        WHERE UserId = ? AND CommentId = ? AND Action = 'D'
        `
		_, err = rp.DB.Exec(query, UserId, CommentId)
		if err != nil {
			return logger.ErrorWrapper("Repository", "LikeComment", "The problem in the LikeComment function (query in Reactions table)", err)
		}
		return nil
	} else if !rp.IsCommentDisliked(UserId, CommentId) && !rp.IsCommentLiked(UserId, CommentId) {
		query := `
        INSERT INTO Commentaries 
        (CommentId, UserId, LikeCount)
        VALUES (?, ?, 1) 
        `
		_, err := rp.DB.Exec(query, CommentId, UserId)
		if err != nil {
			return logger.ErrorWrapper("Repository", "LikeComment", "The problem is to insert like into Commentaries table", err)
		}
		query = `
        INSERT INTO Reactions
        (CommentId, UserId, Action)
        VALUES (?, ?, 'L')
        `
		_, err = rp.DB.Exec(query, CommentId, UserId)
		if err != nil {
			return logger.ErrorWrapper("Repository", "LikeComment", "The problem is to insert into Reactions table", err)
		}
	}
	return nil
}

func (rp *Repository) DislikeComments(UserId, CommentId int) error {
	if rp.IsCommentLiked(UserId, CommentId) {
		query := `
		UPDATE Commentaries
		SET LikeCount = LikeCount - 1
		WHERE UserId = ? and CommentId = ?
		`
		_, err := rp.DB.Exec(query, UserId, CommentId)
		if err != nil {
			return logger.ErrorWrapper("Repository", "DislikeComments", "The problem is to update Commentaries", err)
		}
		query = `
		DELETE FROM Reactions
		WHERE UserId = ? and CommentId = ? and Action = 'L'
		`
		_, err = rp.DB.Exec(query, UserId, CommentId)
		if err != nil {
			return logger.ErrorWrapper("Repository", "DislikeComments", "The problem is delete Commentaries reactions", err)
		}
		return nil
	} else if rp.IsCommentDisliked(UserId, CommentId) {
		query := `
        UPDATE Commentaries
        SET DislikeCount = DislikeCount - 1 
        WHERE CommentId = ? AND UserId = ?
        `
		_, err := rp.DB.Exec(query, CommentId, UserId)
		if err != nil {
			return logger.ErrorWrapper("Repository", "DislikeComments", "The problem in the LikeComment function (query in Commentaries table)", err)
		}
		query = `
        DELETE FROM Reactions
        WHERE UserId = ? AND CommentId = ? AND Action = 'D'
        `
		_, err = rp.DB.Exec(query, UserId, CommentId)
		if err != nil {
			return logger.ErrorWrapper("Repository", "DislikeComments", "The problem in the LikeComment function (query in Reactions table)", err)
		}

		return nil

	} else if !rp.IsCommentLiked(UserId, CommentId) && !rp.IsCommentDisliked(UserId, CommentId) {
		query := `
		INSERT INTO Commentaries
		(UserId, CommentId, DislikeCount)
		VALUES (?, ?, 1)

		`
		_, err := rp.DB.Exec(query, UserId, CommentId)
		if err != nil {
			return logger.ErrorWrapper("Repository", "DislikeComments", "The problem is insert into Commentaries tables", err)
		}
		query = `
		INSERT INTO Reactions
		(UserId, CommentId, Action)
		VALUES (?, ?, 'D')
		`
		_, err = rp.DB.Exec(query, UserId, CommentId)
		if err != nil {
			return logger.ErrorWrapper("Repository", "DislikeComments", "The problem is insert into Reactions tables", err)
		}
	}
	return nil
}

func (rp *Repository) IsCommentLiked(UserId, CommentId int) bool {
	var answer bool
	query := `
	EXISTS (SELECT *
			FROM Reactions
			WHERE UserId = ? and CommentId = ? and Actions = 'L')
	`
	err := rp.DB.QueryRow(query, UserId, CommentId).Scan(&answer)
	if err != nil {
		logger.ErrorWrapper("Repository", "IsCommentLiked", "The problem in the IsCommentLiked function", err)
		return answer
	}
	return answer
}

func (rp *Repository) IsCommentDisliked(UserId, CommentId int) bool {
	var answer bool
	query := `
	EXISTS (SELECT *
			FROM Reactions
			WHERE UserId = ? and CommentId = ? and Actions = 'D')
	`
	err := rp.DB.QueryRow(query, UserId, CommentId).Scan(&answer)
	if err != nil {
		logger.ErrorWrapper("Repository", "IsCommentDisliked", "The problem in the IsCommentDisliked function", err)
		return answer
	}
	return answer
}
