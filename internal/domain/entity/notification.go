package entity

type Notifications struct {
	UserNickname string `json:"userNickname"`
	Action       string `json:"action"`
	PostId       int    `json:"post_id"`
}
