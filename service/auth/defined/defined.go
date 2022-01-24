package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type AccessTokenClaims struct {
	Username   string `json:"user_name"`
	UserId uint `json:"user_id"`
	jwt.StandardClaims
}

type TokenPair struct {
	AccessToken string `json:"accessToken"`
}
