package defined

import (
	"mjo/service/user/defined"
	"time"
)

type InsertRequest struct {
	Name      string   `json:"name" validate:"required"`
	UserName     string   `json:"user_name" validate:"required"`
	Password    string `json:"password" validate:"required"`
	CreatedBy     uint   `json:"created_by"`
}

type UpdateByIdRequest struct {
	Name      string   `json:"name" validate:"required"`
	UserName     string   `json:"user_name" validate:"required"`
	Password    string `json:"password" validate:"required"`
	UpdatedBy     uint   `json:"updated_by"`
}

type DefaultResponse struct {
	Id           uint      `json:"id"`
	Name      string   `json:"name" validate:"required"`
	UserName     string   `json:"user_name" validate:"required"`
	Password    string `json:"password" validate:"required"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      uint `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     uint `json:"updated_by"`
}

func NewDefaultResponse(user *defined.User) *DefaultResponse {
	return &DefaultResponse{
		Id:           user.Id,
		Name:       user.Name.String,
		UserName:        user.UserName.String,
		CreatedAt:      user.CreatedAt,
		CreatedBy:     user.CreatedBy,
		UpdatedAt:      user.UpdatedAt,
		UpdatedBy:     user.UpdatedBy,
	}
}