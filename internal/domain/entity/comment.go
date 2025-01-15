package entity

type Comment struct {
	CommentId int    `json:"comment_id"`
	PostId    int    `json:"post_id"`
	UserId    int    `json:"user_id"`
	Content   string `json:"content"`
	CreateAt  string `json:"create_at"`
}
