package sqllite

import (
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
)

func (rp *Repository) GetAllCategories() ([]entity.Categories, error) {

	statement := "SELECT * FROM Categories"

	rows, err := rp.DB.Query(statement)
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetAllCategories", "Failed to execute query for categories", err)
	}

	var categories []entity.Categories

	for rows.Next() {
		category := entity.Categories{}
		err := rows.Scan(&category.CategoryId, &category.CategoryName)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "GetAllCategories", "Failed to scan category row", err)
		}
		categories = append(categories, category)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetAllCategories", "Error occurred during rows iteration", err)
	}

	if err := rows.Close(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetAllCategories", "Failed to close the row of db", err)
	}

	return categories, nil
}

func (rp *Repository) CreateCategory(categoryName string) error {
	statement := "INSERT INTO Categories (CategoryName) VALUES (?)"

	_, err := rp.DB.Exec(statement, categoryName)
	if err != nil {
		return logger.ErrorWrapper("Repository", "AddCategory", "The problem within the process of creation of the category in db", err)
	}
	return nil

}

func (rp *Repository) DeleteCategory(categoryId int) error {

	statement := `
	PRAGMA foreign_keys = ON;
 	DELETE FROM Categories WHERE CategoryId = ?;
 	`
	_, err := rp.DB.Exec(statement, categoryId)
	if err != nil {
		return logger.ErrorWrapper("Repository", "DeleteCertainPost", "The problem within the process of deleting of the post in db", err)
	}

	return nil
}
