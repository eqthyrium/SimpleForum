package entity

type Reports struct {
	UserId int `json:"user_id"`
	PostId int `json:"post_id"`
}

type ReportInfo struct {
	UserId int    `json:"user_id"`
	Email  string `json:"email"`
	PostId int    `json:"post_id"`
}
