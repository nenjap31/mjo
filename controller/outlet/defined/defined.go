package defined

import (
	"mjo/service/outlet/defined"
	"time"
)

type InsertRequest struct {
	MerchantId      uint   `json:"merchant_id" validate:"required"`
	OutletName     string   `json:"outlet_name" validate:"required"`
	CreatedBy     uint   `json:"created_by"`
}

type UpdateByIdRequest struct {
	MerchantId      uint   `json:"merchant_id" validate:"required"`
	OutletName     string   `json:"outlet_name" validate:"required"`
	UpdatedBy     uint   `json:"updated_by"`
}

type DefaultResponse struct {
	Id           uint      `json:"id"`
	MerchantId      uint   `json:"merchant_id" validate:"required"`
	OutletName     string   `json:"outlet_name" validate:"required"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      uint `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     uint `json:"updated_by"`
}

func NewDefaultResponse(outlet *defined.Outlet) *DefaultResponse {
	return &DefaultResponse{
		Id:           outlet.Id,
		MerchantId:       outlet.MerchantId,
		OutletName:        outlet.OutletName.String,
		CreatedAt:      outlet.CreatedAt,
		CreatedBy:     outlet.CreatedBy,
		UpdatedAt:      outlet.UpdatedAt,
		UpdatedBy:     outlet.UpdatedBy,
	}
}