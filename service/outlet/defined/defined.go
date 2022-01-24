package defined

import (
	"database/sql"
	"time"
)

type Outlet struct {
	Id           uint           `json:"id"`
	MerchantId     uint      `json:"merchant_id"`
	OutletName       sql.NullString           `json:"outlet_name"`
	CreatedAt      time.Time      `json:"created_at"`
	CreatedBy     uint      `json:"created_by"`
	UpdatedAt      time.Time      `json:"updated_at"`
	UpdatedBy     uint      `json:"updated_by"`
}