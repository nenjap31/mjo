package defined

import (
	"time"
)

type LoginRequest struct {
	Username string `json:"user_name" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=1"`
}

type DefaultResponse struct {
	Id                  uint            `json:"id"`
	AccessToken         string            `json:"accessToken"`
	Name               string            `json:"name"`
	UserName            string            `json:"user_name"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type UpdateProfileRequest struct {
	Name   string `json:"name" validate:"required"`
	UserName      string `json:"user_name" validate:"required"`
}
