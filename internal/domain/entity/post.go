package entity

import (
	"database/sql"
	"time"
)

type Posts struct {
	PostId       int            `json:"post_id"`
	UserId       int            `json:"user_id"`
	Title        string         `json:"title"`
	Content      string         `json:"content"`
	Image        sql.NullString `json:"image"`
	LikeCount    int            `json:"like_count"`
	DislikeCount int            `json:"dislike_count"`
	CreatedAt    time.Time      `json:"created_at"`
}
