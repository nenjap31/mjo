package defined

import (
	"database/sql"
	"time"
)

type User struct {
	Id           uint           `json:"id"`
	Name       sql.NullString           `json:"name"`
	UserName      sql.NullString           `json:"user_name"`
	Password     string         `json:"password"`
	CreatedAt      time.Time      `json:"created_at"`
	CreatedBy     uint      `json:"created_by"`
	UpdatedAt      time.Time      `json:"updated_at"`
	UpdatedBy     uint      `json:"updated_by"`
}