package entity

import "time"

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
