package defined

import (
	"mjo/service/transaction/defined"
	"time"
)

type InsertRequest struct {
	MerchantId      uint   `json:"merchant_id" validate:"required"`
	OutletId     uint   `json:"outlet_id" validate:"required"`
	BillTotal     float64   `json:"bill_total" validate:"required"`
	CreatedBy     uint   `json:"created_by"`
}

type UpdateByIdRequest struct {
	MerchantId      uint   `json:"merchant_id" validate:"required"`
	OutletId     uint   `json:"outlet_id" validate:"required"`
	BillTotal     float64   `json:"bill_total" validate:"required"`
	UpdatedBy     uint   `json:"updated_by"`
}

type DefaultResponse struct {
	Id           uint      `json:"id"`
	MerchantId      uint   `json:"merchant_id" validate:"required"`
	OutletId     uint   `json:"outlet_id" validate:"required"`
	BillTotal     float64   `json:"bill_total" validate:"required"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      uint `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     uint `json:"updated_by"`
}

type ReportMonthlyResponse struct {
	MerchantName      string   `json:"merchant_name" validate:"required"`
	Date      string   `json:"date" validate:"required"`
	Omzet     float64   `json:"omzet" validate:"required"`
}

type ReportOutletMonthlyResponse struct {
	MerchantName      string   `json:"merchant_name"`
	OutletName      string   `json:"outlet_name"`
	Date      string   `json:"date"`
	Omzet     float64   `json:"omzet"`
}

type ReportMerchantRequest struct {
	StartDate      string   `json:"start_date" validate:"required"`
	EndDate     string   `json:"end_date" validate:"required"`
	Limit     int   `json:"limit" validate:"required"`
	Offset     int   `json:"offset"`
}

func NewDefaultResponse(transaction *defined.Transaction) *DefaultResponse {
	return &DefaultResponse{
		Id:           transaction.Id,
		MerchantId:       transaction.MerchantId,
		OutletId:        transaction.OutletId,
		BillTotal:        transaction.BillTotal,
		CreatedAt:      transaction.CreatedAt,
		CreatedBy:     transaction.CreatedBy,
		UpdatedAt:      transaction.UpdatedAt,
		UpdatedBy:     transaction.UpdatedBy,
	}
}

func NewReportMonthlyResponse(reportmerchant *defined.ReportMerchant) *ReportMonthlyResponse {
	return &ReportMonthlyResponse{
		MerchantName:           reportmerchant.MerchantName,
		Date:           reportmerchant.Date,
		Omzet:       reportmerchant.Omzet,
	}
}

func NewReportOutletMonthlyResponse(reportoutlet *defined.ReportOutlet) *ReportOutletMonthlyResponse {
	return &ReportOutletMonthlyResponse{
		MerchantName:           reportoutlet.MerchantName,
		OutletName:           reportoutlet.OutletName,
		Date:           reportoutlet.Date,
		Omzet:       reportoutlet.Omzet,
	}
}