package auth

import (
	authApiDefined "mjo/controller/auth/defined"
	//authControllerDefined "mjo/controller/auth/defined"
	authServiceDefined "mjo/service/auth/defined"
	userServiceDefined "mjo/service/user/defined"
	userRepositoryDefined "mjo/repository/user"
	"mjo/repository/util/redis"
	"mjo/config"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Service struct {
	repository     userRepositoryDefined.IRepository
	redis          redis.IRepository
}

type IService interface {
	GenerateToken(username string, password string) (*authApiDefined.DefaultResponse, error)
	RevokeToken(token *jwt.Token) error
	GetProfile(token *jwt.Token) (*userServiceDefined.User, error)
}

func NewService(repository userRepositoryDefined.IRepository, redis redis.IRepository) IService {
	return &Service{
		repository:     repository,
		redis:          redis,
	}
}

func (service *Service) GenerateToken(username string, password string) (*authApiDefined.DefaultResponse, error) {
	user, err := service.repository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		err = errors.New("user not found")
		return nil, err
	}
	var tokenPair *authServiceDefined.TokenPair
	tokenPair, err = GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	err = service.SaveToken(user, tokenPair.AccessToken)
	
	if err != nil {
		return nil, err
	}

	result := authApiDefined.DefaultResponse{
		AccessToken:         tokenPair.AccessToken,
		Id:                  user.Id,
		Name:                user.Name.String,
		UserName:          user.UserName.String,
		CreatedAt:              user.CreatedAt,
		UpdatedAt:              user.UpdatedAt,
	}
	return &result, nil
}

func GenerateTokenPair(user *userServiceDefined.User) (*authServiceDefined.TokenPair, error) {
	claims := authServiceDefined.AccessTokenClaims{
		Username:   user.UserName.String,
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(config.GetConfig().JwtExpired) * time.Minute).Unix(),
			Issuer:    config.GetConfig().Name,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(config.GetConfig().SecretKey))
	if err != nil {
		return nil, err
	}
	return &authServiceDefined.TokenPair{
		AccessToken: accessToken,
	}, nil
}

func (service *Service) SaveToken(user *userServiceDefined.User, token string) error {
	//Check if user token already exist in redis, if so will delete old token
	pattern := user.UserName.String + "|*"
	key, _ := service.redis.Keys(pattern)
	if len(key) > 0 {
		for _, value := range key {
			err := service.redis.Delete(value)
			if err != nil {
				return errors.New("Internal Server Error")
			}
		}
	}

	//Save new token in redis
	r := service.redis.Set(
		user.UserName.String+"|"+token,
		user.UserName.String,
		time.Duration(config.GetConfig().JwtExpired)*time.Minute,
	)
	if r != nil {
		return errors.New("Internal Server Error")
	}

	return nil
}

func (service *Service) RevokeToken(token *jwt.Token) error {
	claims, success := token.Claims.(*authServiceDefined.AccessTokenClaims)
	if !success {
		return errors.New("internal server error")
	}
	jwtTokenRedis := claims.Username + "|" + token.Raw
	err := service.redis.Delete(jwtTokenRedis)
	if err != nil {
		return errors.New("internal server error")
	}
	return nil
}

func (service *Service) GetProfile(token *jwt.Token) (*userServiceDefined.User, error) {
	claims, success := token.Claims.(*authServiceDefined.AccessTokenClaims)
	if !success {
		return nil, errors.New("internal server error")
	}

	data, err := service.repository.FindByUsername(claims.Username)
	if err != nil {
		return nil, err
	}

	return data, nil
}
