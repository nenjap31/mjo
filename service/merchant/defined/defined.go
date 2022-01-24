package defined

import (
	"database/sql"
	"time"
)

type Merchant struct {
	Id           uint           `json:"id"`
	UserId     uint      `json:"user_id"`
	MerchantName       sql.NullString           `json:"merchant_name"`
	CreatedAt      time.Time      `json:"created_at"`
	CreatedBy     uint      `json:"created_by"`
	UpdatedAt      time.Time      `json:"updated_at"`
	UpdatedBy     uint      `json:"updated_by"`
}