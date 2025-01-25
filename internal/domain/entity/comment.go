package entity

type Commentaries struct {
	CommentId    int    `json:"comment_id"`
	PostId       int    `json:"post_id"`
	UserId       int    `json:"user_id"`
	Content      string `json:"content"`
	LikeCount    int    `json:"like_count"`
	DislikeCount int    `json:"dislike_count"`
	CreateAt     string `json:"create_at"`
}
