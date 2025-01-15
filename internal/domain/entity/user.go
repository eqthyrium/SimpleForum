package entity

type User struct {
	UserId   int    `json:"user_id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"` // dlina >= 8 and dlina <=32
	Role     string `json:"role"`
}
