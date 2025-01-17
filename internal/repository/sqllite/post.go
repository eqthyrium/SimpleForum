package sqllite

import (
	"SimpleForum/internal/domain/entity"
	"SimpleForum/pkg/logger"
)

func (rp *Repository) GetAllPosts() ([]entity.Post, error) {
	query := `
    SELECT id, title, content, created_at 
    FROM Posts
    ORDER BY created_at DESC
    `
	rows, err := rp.DB.Query(query)
	if err != nil {
		logger.ErrorWrapper("Repository", "GetAllPosts", "There is problem with query in DB", err)
		return nil, err
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var post entity.Post
		err := rows.Scan(&post.PostId, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			logger.ErrorWrapper("Repository", "GetAllPosts", "", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		logger.ErrorWrapper("Repository", "GetAllPosts", "", err)
		return nil, err
	}

	return posts, nil
}

func (re *Repository) CreatePost(post *entity.Post) error {
	query := `
	INSERT into Posts (UserId, Title, Content, Image) VALUES (?, ?, ?, ?)
	`
	_, err := re.DB.Exec(query, post.UserId, post.Title, post.Content, post.Image)
	if err != nil {
		return logger.ErrorWrapper("Repository", "CreatePost", "The problem within the process of creation post in db", err)
	}
	return nil
}

func (re *Repository) DeletePost(post *entity.Post) error {

}

func (re *Repository) UpdatePost(post *entity.Post) {

}
