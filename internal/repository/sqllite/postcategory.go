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

func (rp *Repository) GetPostIDsOfCertainCategory(categoryId int) ([]int, error) {

	statement := `SELECT PostId FROM PostCategories WHERE CategoryId=?`
	rows, err := rp.DB.Query(statement, categoryId)
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetPostIDsOfCertainCategory", "Failed to get all postIds from the postcategory db", err)
	}

	var postIds []int
	for rows.Next() {
		var postId int
		err = rows.Scan(&postId)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "GetPostIDsOfCertainCategory", "Failed to get all postIds from the postcategory", err)
		}
		postIds = append(postIds, postId)
	}

	if err := rows.Err(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetPostIDsOfCertainCategory", "Error occurred during rows iteration", err)
	}

	if err := rows.Close(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetPostIDsOfCertainCategory", "Failed to close the row of db", err)
	}

	return postIds, nil
}

func (rp *Repository) GetCategoriesOfCertainPost(postId int) ([]int, error) {
	statement := `SELECT CategoryId FROM PostCategories WHERE PostId=?`
	rows, err := rp.DB.Query(statement, postId)
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetPostIDsOfCertainCategory", "Failed to get all postIds from the postcategory db", err)
	}

	var categoryIds []int
	for rows.Next() {
		var categoryId int
		err = rows.Scan(&categoryId)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "GetPostIDsOfCertainCategory", "Failed to get all categoryIds from the postcategory", err)
		}
		categoryIds = append(categoryIds, categoryId)
	}

	if err := rows.Err(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetPostIDsOfCertainCategory", "Error occurred during rows iteration", err)
	}

	if err := rows.Close(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetPostIDsOfCertainCategory", "Failed to close the row of db", err)
	}

	return categoryIds, nil
}
