package domain

import "time"

type User struct {
	UserId   int    `json:"user_id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Post struct {
	PostId       int       `json:"post_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	LikeCount    int       `json:"like_count"`
	DislikeCount int       `json:"dislike_count"`
	CreatedAt    time.Time `json:"created_at"`
	UserId       int       `json:"user_id"`
}
type Comment struct {
	CommentId int    `json:"comment_id"`
	PostId    int    `json:"post_id"`
	UserId    int    `json:"user_id"`
	Content   string `json:"content"`
}

type Category struct {
	CategoryId   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type PostCategory struct {
	PostId     int `json:"post_id"`
	CategoryId int `json:"category_id"`
}

// Notification, Reaction struct must be created
