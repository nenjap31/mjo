package defined

import (
	"time"
)

type Transaction struct {
	Id           uint           `json:"id"`
	MerchantId     uint      `json:"merchant_id"`
	OutletId       uint           `json:"outlet_id"`
	BillTotal       float64           `json:"bill_total"`
	CreatedAt      time.Time      `json:"created_at"`
	CreatedBy     uint      `json:"created_by"`
	UpdatedAt      time.Time      `json:"updated_at"`
	UpdatedBy     uint      `json:"updated_by"`
}

type ReportMerchant struct {
	MerchantName     string      `json:"merchant_name"`
	Date     string      `json:"date"`
	Omzet       float64           `json:"omzet"`
}

type ReportOutlet struct {
	MerchantName     string      `json:"merchant_name"`
	OutletName     string      `json:"outlet_name"`
	Date     string      `json:"date"`
	Omzet       float64           `json:"omzet"`
}