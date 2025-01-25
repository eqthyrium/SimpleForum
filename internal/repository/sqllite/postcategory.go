package sqllite

import "SimpleForum/pkg/logger"

func (rp *Repository) SetPostCategoryRelation(postId, categoryId int) error {
	statement := `INSERT INTO PostCategories (PostId,CategoryId) VALUES (?,?)`
	_, err := rp.DB.Exec(statement, postId, categoryId)
	if err != nil {
		return logger.ErrorWrapper("Repository", "SetPostCategoryRelation", "The problem within the process of setting relation between a post with a category  in db", err)
	}
	return nil
}
