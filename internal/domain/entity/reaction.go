package entity

import (
	"database/sql"
	"time"
)

type Reactions struct {
	UserId    int
	PostId    sql.NullInt32
	CommentId sql.NullInt32
	Action    string
	CreatedAt time.Time
}
