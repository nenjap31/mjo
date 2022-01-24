package defined

import (
	"mjo/service/merchant/defined"
	"time"
)

type InsertRequest struct {
	UserId      uint   `json:"user_id" validate:"required"`
	MerchantName     string   `json:"merchant_name" validate:"required"`
	CreatedBy     uint   `json:"created_by"`
}

type UpdateByIdRequest struct {
	UserId      uint   `json:"user_id" validate:"required"`
	MerchantName     string   `json:"merchant_name" validate:"required"`
	UpdatedBy     uint   `json:"updated_by"`
}

type DefaultResponse struct {
	Id           uint      `json:"id"`
	UserId      uint   `json:"user_id" validate:"required"`
	MerchantName     string   `json:"merchant_name" validate:"required"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      uint `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     uint `json:"updated_by"`
}

func NewDefaultResponse(merchant *defined.Merchant) *DefaultResponse {
	return &DefaultResponse{
		Id:           merchant.Id,
		UserId:       merchant.UserId,
		MerchantName:        merchant.MerchantName.String,
		CreatedAt:      merchant.CreatedAt,
		CreatedBy:     merchant.CreatedBy,
		UpdatedAt:      merchant.UpdatedAt,
		UpdatedBy:     merchant.UpdatedBy,
	}
}