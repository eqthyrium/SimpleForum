package sqllite

import (
	"SimpleForum/internal/domain"
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
)

func (rp *Repository) GetAllCategories() ([]entity.Categories, error) {

	statement := "SELECT * FROM Categories"

	rows, err := rp.DB.Query(statement)
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetAllCategories", "Failed to execute query for categories", err)
	}
	defer rows.Close()

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

	// If no categories found, return an appropriate error
	if len(categories) == 0 {
		return nil, logger.ErrorWrapper("Repository", "GetAllCategories", "No categories found in the database", domain.ErrCategoryNotFound)
	}

	return categories, nil
}
