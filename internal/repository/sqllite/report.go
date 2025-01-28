package sqllite

import (
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
)

func (rp *Repository) GetAllReports() ([]entity.ReportInfo, error) {
	statement := `SELECT u.UserId, u.Email, r.PostId FROM Reports r Inner Join Users u ON r.UserId = u.UserId  WHERE r.PostId != -1`

	rows, err := rp.DB.Query(statement)

	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetAllReports", "Failed to execute query for reports", err)
	}
	reports := make([]entity.ReportInfo, 0)

	for rows.Next() {
		report := entity.ReportInfo{}
		err := rows.Scan(&report.UserId, &report.Email, &report.PostId)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "GetAllReports", "Failed to scan report row", err)
		}
		reports = append(reports, report)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetAllReports", "Error occurred during rows iteration", err)
	}

	if err := rows.Close(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetAllReports", "Failed to close the row of db", err)
	}
	return reports, nil

}
func (rp *Repository) GetAllRequests() ([]entity.ReportInfo, error) {
	statement := `SELECT  u.UserId, u.Email FROM Reports r Inner Join Users u ON r.UserId = u.UserId  WHERE r.PostId = -1`

	rows, err := rp.DB.Query(statement)

	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetAllRequest", "Failed to execute query for request", err)
	}
	requests := make([]entity.ReportInfo, 0)

	for rows.Next() {
		request := entity.ReportInfo{}
		err := rows.Scan(&request.UserId, &request.Email)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "GetAllRequest", "Failed to scan request row", err)
		}
		requests = append(requests, request)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetAllRequest", "Error occurred during rows iteration", err)
	}

	if err := rows.Close(); err != nil {
		return nil, logger.ErrorWrapper("Repository", "GetAllRequest", "Failed to close the row of db", err)
	}
	return requests, nil

}
func (rp *Repository) IsItReportedPost(userId, postId int) (bool, error) {
	var exists bool

	var statement string = `
		SELECT EXISTS(
			SELECT 1 FROM Reports WHERE  UserId = ? AND PostId = ?
		)
	`

	err := rp.DB.QueryRow(statement, userId, postId).Scan(&exists)

	if err != nil {
		return false, logger.ErrorWrapper("Repository", "IsItReportedPost", "Error checking existence of the report", err)
	}

	return exists, nil
}

func (rp *Repository) IsItRequestedToBeModerator(userId int) (bool, error) {

	var exists bool

	var statement string = `
		SELECT EXISTS(
			SELECT 1 FROM Reports WHERE  UserId = ?
		)
	`

	err := rp.DB.QueryRow(statement, userId).Scan(&exists)

	if err != nil {
		return false, logger.ErrorWrapper("Repository", "IsItRequestedToBeModerator", "Error checking existence of the report", err)
	}

	return exists, nil
}

func (rp *Repository) CheckExistenceOfSuchReport(userId, postId int) (bool, error) {

	var exists bool

	var statement string = `
		SELECT EXISTS(
			SELECT 1 FROM Reports WHERE  UserId = ? AND PostId = ?
		)
	`

	err := rp.DB.QueryRow(statement, userId, postId).Scan(&exists)

	if err != nil {
		return false, logger.ErrorWrapper("Repository", "CheckExistenceOfSuchReport", "Error checking existence of the report", err)
	}

	return exists, nil
}

func (rp *Repository) CreateRequestToBeModerator(userId int, falseNumber int) error {

	statement := `
			PRAGMA foreign_keys = OFF;	
			INSERT INTO Reports (UserId, PostId) VALUES(?,?)`
	_, err := rp.DB.Exec(statement, userId, falseNumber)
	if err != nil {
		return logger.ErrorWrapper("Repository", "CreateRequestToBeModerator", "The problem within the process of creating a request of the user to be moderator in db", err)
	}

	return nil
}

func (rp *Repository) CreateReport(userId int, postId int) error {
	statement := `INSERT INTO Reports (UserId, PostId) VALUES(?,?)`
	_, err := rp.DB.Exec(statement, userId, postId)
	if err != nil {
		return logger.ErrorWrapper("Repository", "CreateReport", "The problem within the process of creating a report of the  moderator in db", err)
	}

	return nil
}

func (rp *Repository) DeleteCertainReport(userId, postId int) error {
	statement := `
		PRAGMA foreign_keys = ON;	
	DELETE FROM Reports WHERE UserId = ? AND PostId = ?`
	_, err := rp.DB.Exec(statement, userId, postId)
	if err != nil {
		return logger.ErrorWrapper("Repository", "DeleteCertainReport", "Failed to delete certain report in db", err)
	}
	return nil
}

func (rp *Repository) DeleteRequestToBeModerator(userId, postId int) error {
	statement := `
			PRAGMA foreign_keys = ON;	
		DELETE FROM Reports WHERE UserId = ? AND PostId = ?`
	_, err := rp.DB.Exec(statement, userId, postId)
	if err != nil {
		return logger.ErrorWrapper("Repository", "DeleteRequestToBeModerator", "Failed to delete certain request in db", err)
	}
	return nil
}
