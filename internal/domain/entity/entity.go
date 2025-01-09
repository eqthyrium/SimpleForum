package entity

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
	UserId       int       `json:"user_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Image        string    `json:"image"`
	LikeCount    int       `json:"like_count"`
	DislikeCount int       `json:"dislike_count"`
	CreatedAt    time.Time `json:"created_at"`
}
type Comment struct {
	CommentId    int       `json:"comment_id"`
	PostId       int       `json:"post_id"`
	UserId       int       `json:"user_id"`
	Content      string    `json:"content"`
	LikeCount    int       `json:"like_count"`
	DislikeCount int       `json:"dislike_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type Category struct {
	CategoryId   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type PostCategory struct {
	CategoryId int `json:"category_id"`
	PostId     int `json:"post_id"`
}

type Notification struct {
	UserId        int    `json:"user_id"`
	ResaverUserId int    `json:"resaver_user_id"`
	Content       string `json:"content"`
	CreatedAt     int    `json:"created_at"`
}
